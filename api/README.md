# api


## API requests

##### Add banner to rotation

```bash
curl -X "POST" "http://localhost:7766/banner/add" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
        "bannerId": 1,
        "slotId": 1,
        "description": "banner 1"
      }'
```

Result:

```json
{
  "id": 1,
  "bannerId": 1,
  "slotId": 1,
  "description": "banner 1",
  "createAt": "2019-11-18T19:05:52.023825Z"
}
```
---

##### Set transition for banner

```bash
curl -X "POST" "http://localhost:7766/banner/set-transition" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
        "bannerId": 1,
        "groupId": 1
      }'
```

Result:

```
ok
```
---

##### Selects a banner to display

```bash
curl -X "POST" "http://localhost:7766/banner/select" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
        "slotId": 1,
        "groupId": 1
      }'
```

Result:

```
1
```
---

##### Removes the banner from the rotation

```bash
curl -X "DELETE" "http://localhost:7766/banner/remove/1"
```

Result:

```
ok
```
---