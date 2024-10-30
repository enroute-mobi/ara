package slite

import "encoding/json"

func RewriteValues(jsonPayload []byte) ([]byte, error) {
	jsonMap := make(map[string]interface{})

	err := json.Unmarshal(jsonPayload, &jsonMap)
	if err != nil {
		return nil, err
	}

	rewriteVisit(newRewriteRootParent(), jsonMap)

	return json.Marshal(jsonMap)
}

type rewriteParent interface {
	SetText(text string)
}

type rewriteMapParent struct {
	content map[string]interface{}
	key     string
}

func newRewriteMapParent(content map[string]interface{}, key string) *rewriteMapParent {
	return &rewriteMapParent{content, key}
}

func (parent *rewriteMapParent) SetText(text string) {
	parent.content[parent.key] = text
}

type rewriteRootParent struct {
}

func newRewriteRootParent() *rewriteRootParent {
	return &rewriteRootParent{}
}

func (parent *rewriteRootParent) SetText(text string) {
}

func rewriteVisit(parent rewriteParent, content map[string]interface{}) {
	if len(content) == 1 {
		value, ok := content["value"]
		if ok {
			parent.SetText(value.(string))
			return
		}
	}

	for key, value := range content {
		if mapValue, ok := value.(map[string]interface{}); ok {
			rewriteVisit(newRewriteMapParent(content, key), mapValue)

			if len(mapValue) == 0 {
				delete(content, key)
			}
		} else if arrayValue, ok := value.([]interface{}); ok {
			if len(arrayValue) == 0 {
				delete(content, key)
			} else {
				parent := newRewriteMapParent(content, key)
				for _, entry := range arrayValue {
					switch ent := entry.(type) {
					case string:
						parent.SetText(ent)
					case map[string]interface{}:
						rewriteVisit(parent, ent)
					}
				}
			}
		}
	}
}
