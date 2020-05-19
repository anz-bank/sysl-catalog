package catalog

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/alecthomas/assert"

	"github.com/anz-bank/sysl/pkg/parse"
)

func TestCreateIntegrationDiagramPlantuml(t *testing.T) {
	m, err := parse.NewParser().ParseString(`
App1:
	Endpoint1:
		App2 <- Endpoint2
App2:
	Endpoint2:
		...

`)
	assert.Nil(t, err)

	gen := &Generator{RootModule: m, Fs: afero.NewMemMapFs(), FilesToCreate: map[string]string{}}
	filename := gen.CreateIntegrationDiagramPlantuml(m, "", false)
	contents := gen.FilesToCreate[filename]
	assert.Equal(t,
		"/svg/~1UDgCpa5Bn30G1U1xViNpr5DXgo8UbhBTjegNBKZ5WqW9RMX2aqoPfAY8_rtMYWTFVSVv7ZEJR8v84cpARxLuQflx-bG_5crTeMog6ccAgi6fQL5N3-t5NtNpris_YaE8akFYhD1cK0XHiQBuCIiHUcaLd7n7TdDrUmsjpAYZ29FnisJfq9ERoIiVyIc0e-odaMdnGqcM67UMMDfdRQ8wA_6WU9MZbVqaW8APtjPHoSO5yk9Bl1JpdBr21dGxxFVQZDgUx-Rv3rskbFsZReSqpT5bug3yi3Zx7G00__yzoceA",
		contents)
}
