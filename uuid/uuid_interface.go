package uuid

import (
	"fmt"
	"sync"

	"bitbucket.org/enroute-mobi/ara/config"
	"github.com/satori/uuid"
)

type UUIDGenerator interface {
	NewUUID() string
}

var defaultUUIDGenerator = NewRealUUIDGenerator()

func DefaultUUIDGenerator() UUIDGenerator {
	return defaultUUIDGenerator
}

func SetDefaultUUIDGenerator(generator UUIDGenerator) {
	defaultUUIDGenerator = generator
}

func NewRealUUIDGenerator() UUIDGenerator {
	return &realUUIDGenerator{}
}

type realUUIDGenerator struct{}

func (generator *realUUIDGenerator) NewUUID() string {
	return uuid.NewV4().String()
}

func NewFakeUUIDGenerator() *FakeUUIDGenerator {
	return &FakeUUIDGenerator{realFormat: config.Config.FakeUUIDRealFormat}
}

type FakeUUIDGenerator struct {
	mutex      sync.Mutex
	counter    int
	lastUUID   string
	realFormat bool
}

const fakeRealFormat = "6ba7b814-9dad-11d1-%04x-00c04fd430c8"
const fakeLegacyFormat = "6ba7b814-9dad-11d1-%x-00c04fd430c8"

func (generator *FakeUUIDGenerator) Format() string {
	if generator.realFormat {
		return fakeRealFormat
	} else {
		return fakeLegacyFormat
	}
}

func (generator *FakeUUIDGenerator) NextUUID() string {
	return fmt.Sprintf(generator.Format(), generator.counter)
}

func (generator *FakeUUIDGenerator) NewUUID() string {
	generator.mutex.Lock()
	defer generator.mutex.Unlock()

	uuid := generator.NextUUID()

	generator.lastUUID = uuid
	generator.counter = (generator.counter + 1) % 0xffff

	return uuid
}

func (generator *FakeUUIDGenerator) LastUUID() string {
	return generator.lastUUID
}

type UUIDConsumer struct {
	uuidGenerator UUIDGenerator
}

func NewUUIDConsumer(g UUIDGenerator) UUIDConsumer {
	return UUIDConsumer{uuidGenerator: g}
}

type UUIDInterface interface {
	SetUUIDGenerator(uuidGenerator UUIDGenerator)
	UUIDGenerator() UUIDGenerator
	NewUUID() string
}

func (consumer *UUIDConsumer) SetUUIDGenerator(uuidGenerator UUIDGenerator) {
	consumer.uuidGenerator = uuidGenerator
}

func (consumer *UUIDConsumer) UUIDGenerator() UUIDGenerator {
	if consumer.uuidGenerator == nil {
		consumer.uuidGenerator = DefaultUUIDGenerator()
	}
	return consumer.uuidGenerator
}

func (consumer *UUIDConsumer) NewUUID() string {
	return consumer.UUIDGenerator().NewUUID()
}
