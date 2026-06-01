#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://localhost:8080"
REF="tao"
ADMIN_TOKEN="6ceab96a-8d97-4f2a-8d69-32569a38fc64"
REF_TOKEN="bench-token"
API="$BASE_URL/$REF"
N=80
CREDENTIAL="BENCH:default"
PAYLOAD="/tmp/bench_notify_stop_monitoring.xml"
REF_ID=""

# Helper: POST avec affichage de l'erreur si le code HTTP n'est pas 2xx
api_post() {
    local url="$1"
    local data="$2"
    local auth="${3:-$REF_TOKEN}"
    local response http_code body

    response=$(curl -s -w "\n%{http_code}" -X POST "$url" \
        -H "Content-Type: application/json" \
        -H "Authorization: Token token=$auth" \
        -d "$data")
    http_code=$(echo "$response" | tail -1)
    body=$(echo "$response" | head -n -1)

    if [[ "$http_code" != 2* ]]; then
        echo "ERROR: HTTP $http_code on POST $url" >&2
        echo "Response: $body" >&2
        exit 1
    fi
    echo "$body"
}

# ---------------------------------------------------------------------------
teardown() {
    if [ -n "$REF_ID" ]; then
        echo ""
        echo "--- Teardown: deleting referential '$REF' ($REF_ID) ---"
        curl -sf -X DELETE "$BASE_URL/_referentials/$REF_ID" \
            -H "Authorization: Token token=$ADMIN_TOKEN" >/dev/null
        echo "Done."
    fi
}
trap teardown EXIT

# ---------------------------------------------------------------------------
echo "--- Cleanup: deleting existing referential '$REF' if present ---"
EXISTING_REF_ID=$(curl -sf "$BASE_URL/_referentials" \
    -H "Authorization: Token token=$ADMIN_TOKEN" |
    jq -r ".[] | select(.Slug == \"$REF\") | .Id")
if [ -n "$EXISTING_REF_ID" ] && [ "$EXISTING_REF_ID" != "null" ]; then
    curl -sf -X DELETE "$BASE_URL/_referentials/$EXISTING_REF_ID" \
        -H "Authorization: Token token=$ADMIN_TOKEN" >/dev/null
    echo "Deleted existing referential ($EXISTING_REF_ID)"
fi

# ---------------------------------------------------------------------------
echo "--- Creating referential '$REF' ---"
REF_RESPONSE=$(api_post "$BASE_URL/_referentials" \
    "{\"Slug\":\"$REF\",\"Tokens\":[\"$REF_TOKEN\"]}" \
    "$ADMIN_TOKEN")
REF_ID=$(echo "$REF_RESPONSE" | jq -r '.Id')
echo "Referential ID: $REF_ID"

# ---------------------------------------------------------------------------
echo "--- Creating $N StopAreas ---"
for i in $(seq 1 $N); do
    api_post "$API/stop_areas" \
        "{\"Name\":\"BenchStop $i\",\"Codes\":{\"internal\":\"BENCH:StopPoint:SP:$i:LOC\"}}" >/dev/null
    printf "."
done
echo ""

# ---------------------------------------------------------------------------
echo "--- Creating partner ---"
PARTNER=$(api_post "$API/partners" "{
    \"Slug\": \"bench-partner\",
    \"ConnectorTypes\": [\"siri-stop-monitoring-subscription-collector\"],
    \"Settings\": {
        \"local_credential\": \"$CREDENTIAL\",
        \"remote_code_space\": \"internal\",
        \"remote_url\": \"http://localhost:9999\",
        \"remote_credential\": \"$CREDENTIAL\"
    }
}")
PARTNER_ID=$(echo "$PARTNER" | jq -r '.Id')
echo "Partner ID: $PARTNER_ID"

# ---------------------------------------------------------------------------
echo "--- Creating $N subscriptions (1 per StopArea) and building deliveries ---"
NOW=$(date -u +"%Y-%m-%dT%H:%M:%S.000+00:00")
TODAY=$(date -u +"%Y-%m-%d")

DELIVERIES=""
SUB_COUNT=0
for i in $(seq 1 $N); do
    STOP="BENCH:StopPoint:SP:$i:LOC"
    SUB=$(api_post "$API/partners/$PARTNER_ID/subscriptions" \
        "{\"Kind\":\"StopMonitoringCollect\",\"References\":[{\"Type\":\"StopArea\",\"Code\":{\"internal\":\"$STOP\"}}]}")
    SUB_REF=$(echo "$SUB" | jq -r '.SubscriptionRef')
    SUB_COUNT=$((SUB_COUNT + 1))
    DELIVERIES+="
        <ns2:StopMonitoringDelivery version=\"2.0:FR-IDF-2.4\">
          <ns2:ResponseTimestamp>$NOW</ns2:ResponseTimestamp>
          <ns2:SubscriberRef>BENCH:Subscriber:$i:LOC</ns2:SubscriberRef>
          <ns2:SubscriptionRef>$SUB_REF</ns2:SubscriptionRef>
          <ns2:Status>true</ns2:Status>
          <ns2:MonitoredStopVisit>
            <ns2:RecordedAtTime>$NOW</ns2:RecordedAtTime>
            <ns2:ItemIdentifier>BENCH:Item:$i:LOC</ns2:ItemIdentifier>
            <ns2:MonitoringRef>$STOP</ns2:MonitoringRef>
            <ns2:MonitoredVehicleJourney>
              <ns2:LineRef>BENCH:Line:1:LOC</ns2:LineRef>
              <ns2:DirectionRef>Left</ns2:DirectionRef>
              <ns2:FramedVehicleJourneyRef>
                <ns2:DataFrameRef>$TODAY</ns2:DataFrameRef>
                <ns2:DatedVehicleJourneyRef>BENCH:VehicleJourney:1</ns2:DatedVehicleJourneyRef>
              </ns2:FramedVehicleJourneyRef>
              <ns2:PublishedLineName>Bench Line</ns2:PublishedLineName>
              <ns2:OriginRef>BENCH:StopPoint:SP:1:LOC</ns2:OriginRef>
              <ns2:DestinationRef>BENCH:StopPoint:SP:$N:LOC</ns2:DestinationRef>
              <ns2:Monitored>true</ns2:Monitored>
              <ns2:MonitoredCall>
                <ns2:StopPointRef>$STOP</ns2:StopPointRef>
                <ns2:Order>$i</ns2:Order>
                <ns2:StopPointName>Bench Stop $i</ns2:StopPointName>
                <ns2:VehicleAtStop>false</ns2:VehicleAtStop>
                <ns2:AimedArrivalTime>$NOW</ns2:AimedArrivalTime>
                <ns2:ExpectedArrivalTime>$NOW</ns2:ExpectedArrivalTime>
                <ns2:ArrivalStatus>onTime</ns2:ArrivalStatus>
              </ns2:MonitoredCall>
            </ns2:MonitoredVehicleJourney>
          </ns2:MonitoredStopVisit>
        </ns2:StopMonitoringDelivery>"
    printf "."
done
echo ""
echo "$SUB_COUNT subscriptions created"

# ---------------------------------------------------------------------------
echo "--- Building XML payload ---"

cat >"$PAYLOAD" <<EOF
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns6:NotifyStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo>
        <ns2:ResponseTimestamp>$NOW</ns2:ResponseTimestamp>
        <ns2:ProducerRef>$CREDENTIAL</ns2:ProducerRef>
        <ns2:ResponseMessageIdentifier>BENCH:Message::1:LOC</ns2:ResponseMessageIdentifier>
        <ns2:RequestMessageRef>BENCH:Request::1</ns2:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Notification>
        $DELIVERIES
      </Notification>
      <SiriExtension />
    </ns6:NotifyStopMonitoring>
  </soap:Body>
</soap:Envelope>
EOF

PAYLOAD_BYTES=$(wc -c <"$PAYLOAD")
PAYLOAD_MB=$(awk "BEGIN {printf \"%.3f\", $PAYLOAD_BYTES / 1048576}")
echo "Payload: $PAYLOAD (${PAYLOAD_BYTES} bytes / ${PAYLOAD_MB} MB)"

# ---------------------------------------------------------------------------
echo ""
echo "--- Smoke test ---"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/$REF/siri" \
    -H "Content-Type: application/xml" \
    --data-binary "@$PAYLOAD")
echo "HTTP $HTTP_CODE"
if [ "$HTTP_CODE" != "200" ]; then
    echo "ERROR: expected 200. Check that ara is running (make run)."
    exit 1
fi

# ---------------------------------------------------------------------------
echo ""
echo "--- Benchmark ---"
hey -q 20 -c 20 -z 30s \
    -m POST \
    -H "Content-Type: application/xml" \
    -D "$PAYLOAD" \
    "$BASE_URL/$REF/siri"
