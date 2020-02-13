set -x

package="gitlab.com/sdce/exlib"
proto_package="gitlab.com/sdce/protogo"
mock_folder="mock"
mkdir -p $mock_folder

array=(
       "go.mongodb.org/mongo-driver/mongo Cursor > $mock_folder/mock_cursor.go"
       "github.com/Shopify/sarama SyncProducer > $mock_folder/mock_kafka.go"
       "github.com/Shopify/sarama AsyncProducer > $mock_folder/mock_kafka_async.go"
)

for i in "${array[@]}"; do   # The quotes are necessary here
    eval "mockgen -package mock $i"
done
