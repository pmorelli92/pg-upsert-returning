UPSERT_LOCK=$(echo '{
  method: "POST",
  url: "http://app:8080/upsert-lock",
  body:"{\"id\": \"<ruuid>\"}" | @base64,
  header: {"Content-Type": ["application/json"]}
}' | sed "s/<ruuid>/$(uuidgen | tr '[:upper:]' '[:lower:]')/g")
echo "PGX: EXECUTING FOR LOCKS"
jq -ncM "$UPSERT_LOCK" \
  | vegeta attack -format=json -duration=40s -connections=20 -rate=100 | vegeta encode \
  | vegeta report -type="hist[0,2ms,4ms,6ms,8ms,10ms,15ms]"

UPSERT_CONFLICT=$(echo '{
  method: "POST",
  url: "http://app:8080/upsert-conflict",
  body:"{\"id\": \"<ruuid>\"}" | @base64,
  header: {"Content-Type": ["application/json"]}
}' | sed "s/<ruuid>/$(uuidgen | tr '[:upper:]' '[:lower:]')/g")
echo "PGX: EXECUTING FOR UPSERT CONFLICT"
jq -ncM "$UPSERT_CONFLICT" \
  | vegeta attack -format=json -duration=40s -connections=20 -rate=100 | vegeta encode \
  | vegeta report -type="hist[0,2ms,4ms,6ms,8ms,10ms,15ms]"

UPSERT_NOTHING=$(echo '{
  method: "POST",
  url: "http://app:8080/upsert-donothing",
  body:"{\"id\": \"<ruuid>\"}" | @base64,
  header: {"Content-Type": ["application/json"]}
}' | sed "s/<ruuid>/$(uuidgen | tr '[:upper:]' '[:lower:]')/g")
echo "PGX: EXECUTING FOR UPSERT DO NOTHING"
jq -ncM "$UPSERT_NOTHING" \
 | vegeta attack -format=json -duration=40s -connections=20 -rate=100 | vegeta encode \
 | vegeta report -type="hist[0,2ms,4ms,6ms,8ms,10ms,15ms]"

UPSERT_CTE=$(echo '{
  method: "POST",
  url: "http://app:8080/upsert-cte",
  body:"{\"id\": \"<ruuid>\"}" | @base64,
  header: {"Content-Type": ["application/json"]}
}' | sed "s/<ruuid>/$(uuidgen | tr '[:upper:]' '[:lower:]')/g")
echo "PGX: EXECUTING FOR UPSERT CTE"
jq -ncM "$UPSERT_CTE" \
 | vegeta attack -format=json -duration=40s -connections=20 -rate=100 | vegeta encode \
 | vegeta report -type="hist[0,2ms,4ms,6ms,8ms,10ms,15ms]"

UPSERT_CTE_RANDOM='{
  method: "POST",
  url: "http://app:8080/upsert-cte-random"
}'
echo "PGX: EXECUTING FOR UPSERT CTE RANDOM"
jq -ncM "$UPSERT_CTE_RANDOM" \
 | vegeta attack -format=json -duration=40s -connections=20 -rate=100 | vegeta encode \
 | vegeta report -type="hist[0,2ms,4ms,6ms,8ms,10ms,15ms]"
