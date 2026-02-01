package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	apperrors "github.com/sgaunet/supervisord/internal/errors"
)

// StringExpression replace the python String like "%(var)s" to string.
type StringExpression struct {
	env map[string]string // the environment variable used to replace the var in the python expression
}

// NewStringExpression create a new StringExpression with the environment variables.
func NewStringExpression(envs ...string) *StringExpression {
	se := &StringExpression{env: make(map[string]string)}

	for _, env := range os.Environ() {
		t := strings.SplitN(env, "=", 2)
		se.env["ENV_"+t[0]] = t[1]
	}
	n := len(envs)
	for i := 0; i+1 < n; i += 2 {
		se.env[envs[i]] = envs[i+1]
	}

	hostname, err := os.Hostname()
	if err == nil {
		se.env["host_node_name"] = hostname
	}

	return se
}

// Add adds environment variable (key,value).
func (se *StringExpression) Add(key string, value string) *StringExpression {
	se.env[key] = value
	return se
}

func findVariableEnd(s string, start int, n int) int {
	end := start + 1
	// find variable end indicator
	for end < n && s[end] != ')' {
		end++
	}
	return end
}

func findVariableType(s string, end int, n int) int {
	typ := end + 1
	// find the type of the variable
	for typ < n {
		ch := s[typ]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			break
		}
		typ++
	}
	return typ
}

func (se *StringExpression) substituteVariable(s string, start int, end int, typ int) (string, error) {
	varName := s[start+2 : end]

	varValue, ok := se.env[varName]

	if !ok {
		return "", apperrors.NewEnvVarNotFoundError(varName)
	}
	switch s[typ] {
	case 'd':
		i, err := strconv.Atoi(varValue)
		if err != nil {
			return "", apperrors.NewEnvVarConversionError(varValue)
		}
		return s[0:start] + fmt.Sprintf("%"+s[end+1:typ+1], i) + s[typ+1:], nil
	case 's':
		return s[0:start] + varValue + s[typ+1:], nil
	default:
		return "", apperrors.NewTypeNotImplementedError(string(s[typ]))
	}
}

// Eval substitutes "%(var)s" in given string with evaluated values, and returns resulting string.
func (se *StringExpression) Eval(s string) (string, error) {
	for {
		// find variable start indicator
		start := strings.Index(s, "%(")

		if start == -1 {
			return s, nil
		}

		n := len(s)
		end := findVariableEnd(s, start, n)
		typ := findVariableType(s, end, n)

		// evaluate the variable
		if typ < n {
			var err error
			s, err = se.substituteVariable(s, start, end, typ)
			if err != nil {
				return "", err
			}
		} else {
			return "", apperrors.ErrInvalidStringExpr
		}
	}
}
