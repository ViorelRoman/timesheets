package main

type Config struct {
	Repos struct {
		HG   []string `yaml:"hg"`
		Git  []string `yaml:"git"`
		User string   `yaml:"user"`
	} `yaml:"repos"`
	Calendars struct {
		Calendar      []string `yaml:"calendar"`
		DefaultTicket string   `yaml:"default_ticket"`
	} `yaml:"calendars"`
}

func (c Config) RepoList() []Repo {
	var ret []Repo
	for _, hgRep := range c.Repos.HG {
		ret = append(ret, HGRepo{hgRep})
	}
	for _, gitRep := range c.Repos.Git {
		ret = append(ret, GitRepo{gitRep})
	}
	return ret
}
