package igdb

// Person type
type Person struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	URL         URL         `json:"url"`
	CreatedAt   int         `json:"created_at"` // Unix time in milliseconds
	UpdatedAt   int         `json:"updated_at"` // Unix time in milliseconds
	DOB         int         `json:"dob"`
	Gender      GenderCode  `json:"gender"`
	Country     CountryCode `json:"country"`
	Mugshot     Image       `json:"mug_shot"`
	Bio         string      `json:"bio"`
	Description string      `json:"description"`
	Parent      int         `json:"parent"`
	Homepage    string      `json:"homepage"`
	Twitter     string      `json:"twitter"`
	LinkedIn    string      `json:"linkedin"`
	GooglePlus  string      `json:"google_plus"`
	Facebook    string      `json:"facebook"`
	Instagram   string      `json:"instagram"`
	Tumblr      string      `json:"tumblr"`
	Soundcloud  string      `json:"soundcloud"`
	Pinterest   string      `json:"pinterest"`
	Youtube     string      `json:"youtube"`
	Nicknames   []string    `json:"nicknames"`
	LovesCount  int         `json:"loves_count"`
	Games       []int       `json:"games"`
	Characters  []int       `json:"characters"`
	VoiceActed  []int       `json:"voice_acted"`
}

// GetPerson gets IGDB information for a person identified by its unique IGDB ID.
func (c *Client) GetPerson(id int, opts ...optionFunc) (*Person, error) {
	url, err := c.singleURL(PersonEndpoint, id, opts...)
	if err != nil {
		return nil, err
	}
	var p []Person

	err = c.get(url, &p)
	if err != nil {
		return nil, err
	}

	return &p[0], nil
}

// GetPersons gets IGDB information for a list of people identified by their
// unique IGDB IDs.
func (c *Client) GetPersons(ids []int, opts ...optionFunc) ([]*Person, error) {
	url, err := c.multiURL(PersonEndpoint, ids, opts...)
	if err != nil {
		return nil, err
	}
	var p []*Person

	err = c.get(url, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// SearchPersons searches the IGDB using the given query and returns IGDB information
// for the results. Use functional options for pagination and to sort results by parameter.
func (c *Client) SearchPersons(qry string, opts ...optionFunc) ([]*Person, error) {
	url, err := c.searchURL(PersonEndpoint, qry, opts...)
	if err != nil {
		return nil, err
	}
	var p []*Person

	err = c.get(url, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
