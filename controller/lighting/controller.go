package lighting

import (
	"github.com/reef-pi/reef-pi/controller/connectors"
	"github.com/reef-pi/reef-pi/controller/utils"
	"sync"
	"time"
)

const Bucket = "lightings"

type Config struct {
	DevMode  bool          `json:"dev_mode" yaml:"dev_mode"`
	Interval time.Duration `json:"interval" yaml:"interval"`
}

var DefaultConfig = Config{
	Interval: 15 * time.Second,
}

type Controller struct {
	store     utils.Store
	jacks     *connectors.Jacks
	stopCh    chan struct{}
	telemetry *utils.Telemetry
	config    Config
	running   bool
	mu        *sync.Mutex
	vv        utils.PWM
}

func New(conf Config, jacks *connectors.Jacks, store utils.Store, telemetry *utils.Telemetry) (*Controller, error) {
	var vv utils.PWM
	pwmConf := utils.DefaultPWMConfig
	pwmConf.DevMode = conf.DevMode
	pwm, err := utils.NewPWM(pwmConf)
	if err != nil {
		return nil, err
	}
	vv = pwm
	return &Controller{
		telemetry: telemetry,
		store:     store,
		jacks:     jacks,
		config:    conf,
		stopCh:    make(chan struct{}),
		mu:        &sync.Mutex{},
		vv:        vv,
	}, nil
}

func (c *Controller) Start() {
	go c.StartCycle()
}

func (c *Controller) Stop() {
	c.StopCycle()
}

func (c *Controller) Setup() error {
	if err := c.store.CreateBucket(Bucket); err != nil {
		return err
	}
	lights, err := c.List()
	if err != nil {
		return err
	}
	for _, light := range lights {
		for _, ch := range light.Channels {
			c.telemetry.CreateFeedIfNotExist(light.Name + "-" + ch.Name)
		}
	}
	return nil
}
