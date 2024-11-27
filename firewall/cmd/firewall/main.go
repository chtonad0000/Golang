//go:build !solution

package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Rule struct {
	Endpoint               string   `yaml:"endpoint"`
	ForbiddenUserAgents    []string `yaml:"forbidden_user_agents"`
	ForbiddenHeaders       []string `yaml:"forbidden_headers"`
	RequiredHeaders        []string `yaml:"required_headers"`
	MaxRequestLengthBytes  int      `yaml:"max_request_length_bytes"`
	MaxResponseLengthBytes int      `yaml:"max_response_length_bytes"`
	ForbiddenResponseCodes []int    `yaml:"forbidden_response_codes"`
	ForbiddenRequestRe     []string `yaml:"forbidden_request_re"`
	ForbiddenResponseRe    []string `yaml:"forbidden_response_re"`
}

type Config struct {
	Rules []Rule `yaml:"rules"`
}

func cfgLoad(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func checkRequest(r *http.Request, rule *Rule) (bool, string) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return false, "Error reading request body"
	}
	r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	if rule.MaxRequestLengthBytes > 0 && r.ContentLength > int64(rule.MaxRequestLengthBytes) {
		return false, "Forbidden"
	}

	for _, ua := range rule.ForbiddenUserAgents {
		if matched, _ := regexp.MatchString(ua, r.UserAgent()); matched {
			return false, "Forbidden"
		}
	}
	for _, header := range rule.ForbiddenHeaders {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			if strings.EqualFold(r.Header.Get(strings.TrimSpace(parts[0])), strings.TrimSpace(parts[1])) {
				return false, "Forbidden"
			}
		}
	}
	for _, re := range rule.ForbiddenRequestRe {
		if matched, _ := regexp.MatchString(re, string(bodyBytes)); matched {
			return false, "Forbidden"
		}
	}
	for _, header := range rule.RequiredHeaders {
		if r.Header.Get(header) == "" {
			return false, "Forbidden"
		}
	}

	return true, ""
}

func checkResponse(resp *http.Response, rule *Rule) (bool, string) {
	bodyBytes, _ := io.ReadAll(resp.Body)
	for _, code := range rule.ForbiddenResponseCodes {
		if resp.StatusCode == code {
			return false, "Forbidden"
		}
	}

	for _, re := range rule.ForbiddenResponseRe {
		if matched, _ := regexp.MatchString(re, string(bodyBytes)); matched {
			return false, "Forbidden"
		}
	}

	if rule.MaxResponseLengthBytes > 0 && int(resp.ContentLength) > rule.MaxResponseLengthBytes {
		return false, "Forbidden"
	}

	resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	return true, ""
}

func main() {
	serviceAddr := flag.String("service-addr", "", "Address of the service to protect")
	firewallAddr := flag.String("addr", "", "Firewall address")
	confPath := flag.String("conf", "", "Path to the YAML config file")
	flag.Parse()
	var config *Config
	if *confPath != "" {
		var err error
		config, err = cfgLoad(*confPath)
		if err != nil {
			log.Fatal("Error loading config: ", err)
		}
	} else {
		config = &Config{}
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, rule := range config.Rules {
			if r.URL.Path == rule.Endpoint {
				pass, message := checkRequest(r, &rule)
				if !pass {
					http.Error(w, message, http.StatusForbidden)
					return
				}
			}
		}

		req, err := http.NewRequest(r.Method, *serviceAddr+r.URL.Path, r.Body)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}
		req.Header = r.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Error proxying request", http.StatusInternalServerError)
			return
		}
		defer func(Body io.ReadCloser) {
			errClose := Body.Close()
			if errClose != nil {
				return
			}
		}(resp.Body)
		for _, rule := range config.Rules {
			if r.URL.Path == rule.Endpoint {
				pass, message := checkResponse(resp, &rule)
				if !pass {
					http.Error(w, message, http.StatusForbidden)
					return
				}
			}
		}

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(body)
		if err != nil {
			http.Error(w, "Error writing response body", http.StatusInternalServerError)
		}
	})

	err := http.ListenAndServe(*firewallAddr, nil)
	if err != nil {
		return
	}
}
