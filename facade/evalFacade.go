package facade

// ArgsEvaluationFacade represents the DTO struct used in the auth facade constructor
type ArgsEvaluationFacade struct {
	ApiInterface string
	PprofEnabled bool
}

type evaluationFacade struct {
	apiInterface string
	pprofEnabled bool
}

// NewEvaluationFacade returns a new instance of authFacade
func NewEvaluationFacade(args ArgsEvaluationFacade) (*evaluationFacade, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &evaluationFacade{
		apiInterface: args.ApiInterface,
		pprofEnabled: args.PprofEnabled,
	}, nil
}

// checkArgs check the arguments of an ArgsNewWebServer
func checkArgs(args ArgsEvaluationFacade) error {

	return nil
}

// RestApiInterface returns the interface on which the rest API should start on, based on the flags provided.
// The API will start on the DefaultRestInterface value unless a correct value is passed or
//  the value is explicitly set to off, in which case it will not start at all
func (af *evaluationFacade) RestApiInterface() string {
	return af.apiInterface
}

// PprofEnabled returns if profiling mode should be active or not on the application
func (af *evaluationFacade) PprofEnabled() bool {
	return af.pprofEnabled
}

// IsInterfaceNil returns true if there is no value under the interface
func (af *evaluationFacade) IsInterfaceNil() bool {
	return af == nil
}
