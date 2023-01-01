package proc

type InstallParams struct {
	ConfigPath string
	IsUser bool
}

type InstallClient struct {

}

func NewInstallClient() InstallClient {
	return InstallClient{}
}

func (c InstallClient) Run(params InstallParams) error {

	return nil
}