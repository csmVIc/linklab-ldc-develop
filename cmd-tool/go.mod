module linklab/device-control-v2/cmd-tool

go 1.15

require (
	github.com/docker/cli v20.10.7+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
