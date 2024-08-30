# Custom mocks for the pandora repository
mocks :
	mockery --name='(.*)' --case=underscore --dir=redis --output=redis/mocks
	mockery --name='(.*)' --case=underscore --dir=httpclient --output=httpclient/mocks
	mockery --name='(.*)' --case=underscore --dir=sdks/onfido --output=sdks/onfido/mocks
	mockery --name='(.*)' --case=underscore --dir=sdks/digilocker --output=sdks/digilocker/mocks
	mockery --name='(.*)' --case=underscore --dir=sdks/dms --output=sdks/dms/mocks
	mockery --name='(.*)' --case=underscore --dir=sdks/deposit --output=sdks/deposit/mocks
	mockery --name='(.*)' --case=underscore --dir=sdks/profile --output=sdks/profile/mocks
	mockery --name='(.*)' --case=underscore --dir=sdks/idfy --output=sdks/idfy/mocks