set -e

TEST_RESULT_DIR="${TEST_RESULTS:-./tmp}"
mkdir -p ${TEST_RESULT_DIR}

PKG_LIST=$(go list ./...\
  | grep -v\
      -e /mock\
      -e /cmd\
      -e /gen\
      -e /testdata\
      -e /testenv\
      -e /handlers/middleware\
      -e /util/test\
      -e /util/logger\
      -e /util/rollbar\
      -e /util/gorm\
      -e /util/kafkacg\
      -e /util/httpclient\
  | tr '\n' ',')

echo "---------------------------------------------------------------"
echo "Test:"
echo "---------------------------------------------------------------"
go test -race\
  -coverpkg=${PKG_LIST}\
  -coverprofile ${TEST_RESULT_DIR}/.testCoverage.txt\
  ./...

echo "---------------------------------------------------------------"
echo "Total coverage:"
echo "---------------------------------------------------------------"
go tool cover -func ${TEST_RESULT_DIR}/.testCoverage.txt | grep total | awk '{print $3}'
go tool cover -html=${TEST_RESULT_DIR}/.testCoverage.txt -o ${TEST_RESULT_DIR}/coverage.html
