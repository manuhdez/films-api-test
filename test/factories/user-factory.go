package factories

import (
	"github.com/google/uuid"
	"syreclabs.com/go/faker"

	"github.com/manuhdez/films-api-test/internal/domain/user"
)

func User() user.User {
	return user.New(
		uuid.New(),
		faker.Internet().UserName(),
		faker.Internet().Password(6, 12),
	)
}
