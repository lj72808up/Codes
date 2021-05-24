package models

type RobotAns struct {
	Type     int    `json:"type"`
	AnswerId string `json:"answerId"`
	Data     string `json:"data"`
}

type RobotQuery struct {
	Type     int      `json:"type"`
	AnswerId string   `json:"answerId"`
	Data     []string `json:"data"`
}

type XiaopMsg struct {
	Pubid   string `json:"pubid"`
	Token   string `json:"token"`
	Ts      int64  `json:"ts"`
	To      string `json:"to"`
	Content string `json:"content"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Url     string `json:"url"`
	Tp      string `json:"tp"`
}

type XiaopApply struct {
	Uid       string `json:"uid" orm:"pk"`
	Boss      string `json:"boss"`
	Apartment string `json:"apartment"`
	Position  string `json:"position"`
	Name      string `json:"name"`
	Status    int    `json:"status"`
	Ts        int64  `json:"ts"`
	Reason    string `json:"reason"`
	Detail    string `json:"detail"`
	Role      string `json:"role"`
}
