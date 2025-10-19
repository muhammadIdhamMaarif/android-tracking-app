> Standard library imports needed:

```go
import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"strconv"
)
```

### Run it

```bash
go mod init example.com/locserver
go mod tidy
go run main.go
```

Server listens on `0.0.0.0:5000`.

### Quick test

```bash
curl -X POST http://localhost:5000/api/loc \
  -H "Content-Type: application/json" \
  -d '{"timestamp": 1697040000000, "device_id":"phone1", "lat":-7.982, "lon":112.63, "accuracy":8.5, "speed":0.3}'
```

You should see `{"status":"ok"}` and a new row in `locations.csv`.

