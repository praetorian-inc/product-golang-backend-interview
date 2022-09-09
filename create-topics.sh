docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic ingest

docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic domainEvent

docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic subdomainEvent