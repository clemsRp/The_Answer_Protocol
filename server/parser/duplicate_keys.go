package parser 

import (
	"errors"
	"bytes"
	"encoding/json"
	"io"
	"fmt"
)

func IsValidKeys(file []byte) error {
    // Unmarshal into generic structure
    var world []map[string]any
    if err := json.Unmarshal(file, &world); err != nil || len(world) == 0 {
        return errors.New("Invalid file: JSON file must be parsable")
    }

    // Initialize the decoder
    dec := json.NewDecoder(bytes.NewReader(file))

    var stack []map[string]bool

    for {
        // Read next token
        t, err := dec.Token()
        if err == io.EOF {
            break
        }
        if err != nil {
            return errors.New("Invalid file: JSON file must be parsable")
        }

        // Handle delimiters
        if delim, ok := t.(json.Delim); ok {
            if delim == '{' {
                stack = append(stack, make(map[string]bool))

            } else if delim == '}' {
                if len(stack) > 0 {
                    stack = stack[:len(stack)-1]
                }
            }
            continue
        }

        // Handle keys
        if str, ok := t.(string); ok && len(stack) > 0 {
            currentLevelKeys := stack[len(stack)-1]

            // Compare counts to detect duplicate keys
            if currentLevelKeys[str] {
                return fmt.Errorf("Invalid map: duplicate key '%s' detected", str)
            }
            currentLevelKeys[str] = true

            // Read next token
            next, err := dec.Token()
            if err != nil {
                return errors.New("Invalid file: JSON file must be parsable")
            }

            if d, ok := next.(json.Delim); ok && d == '{' {
                stack = append(stack, make(map[string]bool))
            } else {
                // Skip the value to avoid sub keys or lists
                if err := skipNextValue(dec, next); err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

func skipNextValue(dec *json.Decoder, currentToken json.Token) error {
    // Check if token is a single value
    _, ok := currentToken.(json.Delim)
    if !ok {
        return nil
    }

    // Skip object or list until matching closure
    depth := 1
    for depth > 0 {
        t, err := dec.Token()
        if err != nil {
            return err
        }

        if delim, ok := t.(json.Delim); ok {
            if delim == '{' || delim == '[' {
                depth++

            } else if delim == '}' || delim == ']' {
                depth--
            }
        }
    }
    return nil
}

func skipValue(dec *json.Decoder) error {
	// Read next token
	t, err := dec.Token()
	if err != nil {
		return err
	}

	// Check if token is a single value
	_, ok := t.(json.Delim)
	if !ok {
		return nil
	}

	// Skip object or list until matching closure
	depth := 1
	for depth > 0 {
		t, err = dec.Token()
		if err != nil {
			return err
		}

		if d, ok := t.(json.Delim); ok {
			if d == '{' || d == '[' {
				depth++

			} else if d == '}' || d == ']' {
				depth--
			}
		}
	}
	return nil
}
