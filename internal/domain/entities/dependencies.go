package entities

var (
	CompanyRepositorySvc CompanyRepository
)

func Setup(
	companyRepository CompanyRepository,
) {
	CompanyRepositorySvc = companyRepository
}
