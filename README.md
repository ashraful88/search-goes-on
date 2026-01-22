# Search API - Elasticsearch Query Builder

A production-ready search API built with Go and Elasticsearch, featuring dynamic query generation from YAML configuration and support for full-text search with faceted filtering.

## Features

- **Full-text search** across multiple fields with relevance scoring
- **Faceted filtering** with category, brand, price range, rating, and stock filters
- **YAML-based configuration** for search fields and filters
- **Aggregations** for faceted search results (counts, ranges, averages)
- **Pagination & sorting** with customizable limits and sort fields
- **Query context & filter context** separation for optimal Elasticsearch performance
- **JSON API compliant** responses
- **Structured logging** with detailed request tracking
- **Docker support** with single-node and multi-node Elasticsearch clusters

## Architecture

```
cmd/search-api/
├── main.go                          # Application entry point
└── internal/api/
    ├── router.go                    # Route definitions
    ├── request_handlers.go          # HTTP handlers
    ├── response_helper.go           # JSON API formatting
    └── logger.go                    # Structured logging

internal/search/
├── query.go                         # Query builder with YAML config
├── search.go                        # Elasticsearch query execution
├── filters.go                       # Filter configuration management
└── config/
    └── search_params.yaml           # Search & filter configuration

internal/platform/searchengine/
└── elasticsearch.go                 # Elasticsearch client setup
```

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Elasticsearch 7.x

### 1. Start Elasticsearch

Single-node cluster:
```bash
docker-compose -f docker-stack/docker-compose-single-es.yml up -d
```

Multi-node cluster:
```bash
docker-compose -f docker-stack/docker-compose.yml up -d
```

### 2. Configure Environment

Copy the example environment file:
```bash
cp example.env .env
```

Update the following variables in `.env`:
```bash
MAIN_INDEX=products
ELASTICSEARCH_URL=http://localhost:9200
FILTER_FILE_PATH=./internal/search/config
```

### 3. Run the Application

```bash
go run cmd/search-api/main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### 1. Health Check

```bash
GET /v1/health
```

**Response:**
```json
{
  "status": "ok"
}
```

### 2. Get Available Filters

```bash
GET /v1/filters
```

**Response:**
```json
{
  "data": {
    "type": "filters",
    "id": "available-filters",
    "attributes": {
      "filters": [
        {
          "name": "category",
          "label": "Category",
          "type": "terms",
          "field": "category.keyword",
          "multi": true
        },
        {
          "name": "brand",
          "label": "Brand",
          "type": "terms",
          "field": "brand.keyword",
          "multi": true
        },
        {
          "name": "price",
          "label": "Price",
          "type": "range",
          "field": "price",
          "options": [
            {"key": "0-50", "from": 0, "to": 50},
            {"key": "50-100", "from": 50, "to": 100},
            {"key": "100-200", "from": 100, "to": 200},
            {"key": "200-*", "from": 200}
          ]
        }
      ]
    }
  },
  "meta": {
    "count": 5
  }
}
```

### 3. Search Products

```bash
GET /v1/search?q={query}&category={cat}&brand={brand}&price={range}&limit={n}&from={offset}&sort={field}&order={asc|desc}
```

**Query Parameters:**

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `q` | string | Full-text search query | `q=laptop` |
| `category` | string[] | Filter by category (multi-select) | `category=electronics&category=computers` |
| `brand` | string[] | Filter by brand (multi-select) | `brand=apple&brand=samsung` |
| `price` | string | Filter by price range | `price=100-200` |
| `rating` | string | Filter by rating | `rating=4+` |
| `in_stock` | boolean | Filter by stock availability | `in_stock=true` |
| `limit` | integer | Number of results (default: 40) | `limit=20` |
| `from` | integer | Offset for pagination (default: 0) | `from=20` |
| `sort` | string | Sort field | `sort=price` |
| `order` | string | Sort order (asc/desc) | `order=asc` |

**Example Requests:**

1. Simple text search:
```bash
curl "http://localhost:8080/v1/search?q=laptop"
```

2. Search with filters:
```bash
curl "http://localhost:8080/v1/search?q=laptop&category=electronics&price=100-200&brand=apple"
```

3. Pagination and sorting:
```bash
curl "http://localhost:8080/v1/search?q=laptop&limit=20&from=0&sort=price&order=asc"
```

**Response Structure:**
```json
{
  "hits": {
    "total": {
      "value": 150,
      "relation": "eq"
    },
    "hits": [
      {
        "_index": "products",
        "_id": "1",
        "_score": 4.5,
        "_source": {
          "name": "MacBook Pro",
          "description": "Powerful laptop...",
          "price": 1299,
          "category": "electronics",
          "brand": "apple",
          "rating": 4.8
        }
      }
    ]
  },
  "aggregations": {
    "category": {
      "buckets": [
        {"key": "electronics", "doc_count": 120},
        {"key": "computers", "doc_count": 80}
      ]
    },
    "price_ranges": {
      "buckets": [
        {"key": "0-50", "doc_count": 25},
        {"key": "100-200", "doc_count": 45}
      ]
    }
  }
}
```

## Configuration

### Search Fields Configuration

Edit [internal/search/config/search_params.yaml](internal/search/config/search_params.yaml):

```yaml
search:
  fields:
    - name
    - description
    - brand
    - category
  default_operator: OR
  minimum_should_match: 1
```

- **fields**: List of fields to search when a query string is provided
- **default_operator**: OR (any term matches) or AND (all terms must match)
- **minimum_should_match**: Minimum number of terms that must match

### Filter Configuration

```yaml
filters:
  category:
    label: Category
    field: category.keyword
    type: terms
    multi: true
```

**Filter Types:**

1. **terms**: Multi-value exact match (e.g., categories, brands)
   - `multi: true` allows multiple values
   - Uses Elasticsearch `terms` query

2. **term**: Single-value exact match (e.g., in stock boolean)
   - `multi: false` allows only one value
   - Uses Elasticsearch `term` query

3. **range**: Numeric or date range filters (e.g., price, rating)
   - Requires `ranges` configuration with `from` and `to` values
   - Uses Elasticsearch `range` query

### Aggregations Configuration

```yaml
aggregations:
  category:
    type: terms
    field: category.keyword
    size: 50
```

**Aggregation Types:**

- **terms**: Count documents by distinct values
- **range**: Count documents in numeric ranges
- **avg**: Calculate average of numeric field

## How It Works

### Query Building Process

1. **Query Context (affects scoring)**:
   - The `q` parameter triggers a `multi_match` query
   - Searches across configured fields (name, description, brand, category)
   - Results ranked by relevance score

2. **Filter Context (doesn't affect scoring)**:
   - URL parameters matching filter names apply filters
   - Multiple values for same filter treated as OR (e.g., `category=a&category=b`)
   - Different filters combined with AND
   - Filters don't impact relevance scoring (faster execution)

3. **Example Generated Query**:

For request: `GET /v1/search?q=laptop&category=electronics&price=100-200`

```json
{
  "query": {
    "bool": {
      "must": [
        {
          "multi_match": {
            "query": "laptop",
            "fields": ["name", "description", "brand", "category"],
            "type": "best_fields",
            "operator": "OR"
          }
        }
      ],
      "filter": [
        {
          "terms": {
            "category.keyword": ["electronics"]
          }
        },
        {
          "range": {
            "price": {
              "gte": 100,
              "lt": 200
            }
          }
        }
      ]
    }
  },
  "size": 40,
  "from": 0,
  "aggs": {
    "category": {
      "terms": {"field": "category.keyword", "size": 50}
    },
    "price_ranges": {
      "range": {
        "field": "price",
        "ranges": [
          {"key": "0-50", "from": 0, "to": 50},
          {"key": "100-200", "from": 100, "to": 200}
        ]
      }
    }
  }
}
```

## Development

### Project Structure

- **cmd/search-api**: Application entry point and API layer
- **internal/search**: Core search logic and query building
- **internal/platform/searchengine**: Elasticsearch client wrapper
- **docker-stack**: Docker Compose configurations for local development

### Adding New Filters

1. Update [search_params.yaml](internal/search/config/search_params.yaml):
```yaml
filters:
  new_filter:
    label: New Filter
    field: new_field.keyword
    type: terms
    multi: true
```

2. Add corresponding aggregation if needed:
```yaml
aggregations:
  new_filter:
    type: terms
    field: new_field.keyword
    size: 20
```

3. Restart the application - filters are loaded from YAML on startup

### Running Tests

```bash
go test ./...
```

## Elasticsearch Index Setup

Your Elasticsearch index should have a mapping similar to:

```json
{
  "mappings": {
    "properties": {
      "name": {"type": "text"},
      "description": {"type": "text"},
      "category": {
        "type": "text",
        "fields": {
          "keyword": {"type": "keyword"}
        }
      },
      "brand": {
        "type": "text",
        "fields": {
          "keyword": {"type": "keyword"}
        }
      },
      "price": {"type": "float"},
      "rating": {"type": "float"},
      "in_stock": {"type": "boolean"},
      "created_at": {"type": "date"}
    }
  }
}
```

## Performance Tips

1. **Use filter context**: Filters don't affect scoring and are cached by Elasticsearch
2. **Keyword fields**: Use `.keyword` fields for exact matching and aggregations
3. **Limit aggregation size**: Set appropriate `size` in aggregations config
4. **Pagination**: Use `from` and `limit` for large result sets
5. **Index optimization**: Ensure proper mapping and analyzers in Elasticsearch

## Troubleshooting

### Connection Issues

If you get connection errors to Elasticsearch:
1. Check Elasticsearch is running: `curl http://localhost:9200`
2. Verify `ELASTICSEARCH_URL` in `.env`
3. Check Docker containers: `docker-compose ps`

### Query Not Working

1. Check logs for the generated Elasticsearch query
2. Test the query directly in Kibana Dev Tools
3. Verify field names match your index mapping
4. Check filter configuration in YAML

### No Results

1. Ensure your index has data
2. Check if filters are too restrictive
3. Try a simple query without filters: `GET /v1/search?q=*`

## License

MIT

## Contributing

Contributions are welcome! Please submit pull requests or open issues for bugs and feature requests.
