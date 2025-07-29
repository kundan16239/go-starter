#!/bin/bash

# Usage: ./scripts/generate_domain.sh <domain>
# Example: ./scripts/generate_domain.sh order

DOMAIN=$1
if [ -z "$DOMAIN" ]; then
  echo "Usage: $0 <domain>"
  exit 1
fi

# Capitalize first letter for struct name
CAP_DOMAIN="$(tr '[:lower:]' '[:upper:]' <<< ${DOMAIN:0:1})${DOMAIN:1}"

# Package name must be lowercase
PKG_DOMAIN="$(echo "$DOMAIN" | tr '[:upper:]' '[:lower:]')"

# UpperCamelCase for constructor naming (e.g., userComment => UserComment)
PASCAL_CASE_DOMAIN="$(tr '[:lower:]' '[:upper:]' <<< ${DOMAIN:0:1})${DOMAIN:1}"

mkdir -p internal/$DOMAIN

# 1. model.go
cat > internal/$DOMAIN/${DOMAIN}Model.go <<EOF
package $DOMAIN

import "go.mongodb.org/mongo-driver/bson/primitive"

type ${CAP_DOMAIN} struct {
    ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    // Add more fields here
}

func (o *${CAP_DOMAIN}) ToResponse() ${CAP_DOMAIN}Response {
    return ${CAP_DOMAIN}Response{
        ID: o.ID.Hex(),
        // Map more fields here
    }
}
EOF

# 2. dto.go
cat > internal/$DOMAIN/${DOMAIN}Dto.go <<EOF
package $DOMAIN

type Create${CAP_DOMAIN}Request struct {
    // Add fields for creation
}

type Update${CAP_DOMAIN}Request struct {
    // Add fields for update
}

type ${CAP_DOMAIN}Response struct {
    ID string `json:"id"`
    // Add more fields here
}

type ${CAP_DOMAIN}ListResponse struct {
    ${CAP_DOMAIN}s []${CAP_DOMAIN}Response `json:"${DOMAIN}s"`
    Count  int64           `json:"count"`
}
EOF

# 3. repository.go (implementation)
cat > internal/$DOMAIN/${DOMAIN}Repository.go <<EOF
package $DOMAIN

import (
    "context"
    "errors"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "project-structure/pkg/shared/repository"
)

type Repository interface {
    Create(ctx context.Context, entity *${CAP_DOMAIN}) error
    GetByID(ctx context.Context, id primitive.ObjectID) (*${CAP_DOMAIN}, error)
    Update(ctx context.Context, entity *${CAP_DOMAIN}) error
    Delete(ctx context.Context, id primitive.ObjectID) error
    List(ctx context.Context, skip, limit int64, findOptions *options.FindOptions) ([]*${CAP_DOMAIN}, int64, error)
}

type MongoRepository struct {
    *repository.BaseRepository
}

func New${PASCAL_CASE_DOMAIN}Repository(db *mongo.Database) Repository {
    return &MongoRepository{
        BaseRepository: repository.NewBaseRepository(db, "${DOMAIN}s"),
    }
}

func (r *MongoRepository) Create(ctx context.Context, entity *${CAP_DOMAIN}) error {
    id, err := r.BaseRepository.Create(ctx, entity)
    if err != nil {
        return fmt.Errorf("failed to create $DOMAIN: %w", err)
    }
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("failed to parse object ID: %w", err)
    }
    entity.ID = objectID
    return nil
}

func (r *MongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*${CAP_DOMAIN}, error) {
    filter := bson.M{"_id": id}
    doc, err := r.BaseRepository.GetOne(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to get $DOMAIN by ID: %w", err)
    }
    if doc == nil {
        return nil, errors.New("$DOMAIN not found")
    }
    var entity ${CAP_DOMAIN}
    err = repository.DocumentToStruct(doc, &entity)
    if err != nil {
        return nil, fmt.Errorf("failed to convert document to $DOMAIN: %w", err)
    }
    return &entity, nil
}

func (r *MongoRepository) Update(ctx context.Context, entity *${CAP_DOMAIN}) error {
    filter := bson.M{"_id": entity.ID}
    update := bson.M{"$set": entity}
    modifiedCount, err := r.BaseRepository.UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to update $DOMAIN: %w", err)
    }
    if modifiedCount == 0 {
        return errors.New("$DOMAIN not found")
    }
    return nil
}

func (r *MongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
    filter := bson.M{"_id": id}
    deletedCount, err := r.BaseRepository.DeleteOne(ctx, filter)
    if err != nil {
        return fmt.Errorf("failed to delete $DOMAIN: %w", err)
    }
    if deletedCount == 0 {
        return errors.New("$DOMAIN not found")
    }
    return nil
}

func (r *MongoRepository) List(ctx context.Context, skip, limit int64, findOptions *options.FindOptions) ([]*${CAP_DOMAIN}, int64, error) {
    total, err := r.BaseRepository.CountDocuments(ctx, bson.M{})
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count ${DOMAIN}s: %w", err)
    }
    if findOptions == nil {
        findOptions = options.Find()
    }
    findOptions.SetSkip(skip)
    findOptions.SetLimit(limit)
    documents, err := r.BaseRepository.GetAll(ctx, bson.M{}, findOptions)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get ${DOMAIN}s: %w", err)
    }
    var entities []*${CAP_DOMAIN}
    err = repository.DocumentsToStructs(documents, &entities)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to convert documents to ${DOMAIN}s: %w", err)
    }
    return entities, total, nil
}
EOF

# 4. service.go
cat > internal/$DOMAIN/${DOMAIN}Service.go <<EOF
package $DOMAIN

import (
    "context"
    "errors"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
    repo Repository
}

func New${PASCAL_CASE_DOMAIN}Service(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Create${CAP_DOMAIN}(ctx context.Context, req Create${CAP_DOMAIN}Request) (*${CAP_DOMAIN}Response, error) {
    entity := &${CAP_DOMAIN}{
        // Map fields from req
    }
    if err := s.repo.Create(ctx, entity); err != nil {
        return nil, err
    }
    resp := entity.ToResponse()
    return &resp, nil
}

func (s *Service) Get${CAP_DOMAIN}ByID(ctx context.Context, id string) (*${CAP_DOMAIN}Response, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid $DOMAIN ID")
    }
    entity, err := s.repo.GetByID(ctx, oid)
    if err != nil {
        return nil, err
    }
    resp := entity.ToResponse()
    return &resp, nil
}

func (s *Service) Update${CAP_DOMAIN}(ctx context.Context, id string, req Update${CAP_DOMAIN}Request) (*${CAP_DOMAIN}Response, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid $DOMAIN ID")
    }
    entity, err := s.repo.GetByID(ctx, oid)
    if err != nil {
        return nil, err
    }
    // Map fields from req to entity
    if err := s.repo.Update(ctx, entity); err != nil {
        return nil, err
    }
    resp := entity.ToResponse()
    return &resp, nil
}

func (s *Service) Delete${CAP_DOMAIN}(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return errors.New("invalid $DOMAIN ID")
    }
    return s.repo.Delete(ctx, oid)
}

func (s *Service) List${CAP_DOMAIN}s(ctx context.Context, skip, limit int64) (*${CAP_DOMAIN}ListResponse, error) {
    if limit < 1 || limit > 100 {
        limit = 10
    }
    findOptions := options.Find().SetSkip(skip).SetLimit(limit)
    entities, total, err := s.repo.List(ctx, skip, limit, findOptions)
    if err != nil {
        return nil, err
    }
    resp := make([]${CAP_DOMAIN}Response, len(entities))
    for i, o := range entities {
        resp[i] = o.ToResponse()
    }
    return &${CAP_DOMAIN}ListResponse{
        ${CAP_DOMAIN}s: resp,
        Count:  total,
    }, nil
}
EOF

# 5. handler.go
cat > internal/$DOMAIN/${DOMAIN}Handler.go <<EOF
package $DOMAIN

import (
	"net/http"
	"strconv"

	sharedErr "project-structure/pkg/shared/errors"
	"project-structure/pkg/shared/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	validate *validator.Validate
	logger   logger.Logger
	service  *Service
}

func New${PASCAL_CASE_DOMAIN}Handler(validate *validator.Validate, logger logger.Logger, service *Service) *Handler {
	return &Handler{validate: validate, logger: logger, service: service}
}

func (h *Handler) Create${CAP_DOMAIN}(c *gin.Context) {
	var req Create${CAP_DOMAIN}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid request body")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid request body")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	entity, err := h.service.Create${CAP_DOMAIN}(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create $DOMAIN", "error", err)
		appErr := sharedErr.NewDatabaseError("Failed to create $DOMAIN")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "${CAP_DOMAIN} created successfully", "data": entity})
}

func (h *Handler) Get${CAP_DOMAIN}(c *gin.Context) {
	id := c.Param("id")
	entity, err := h.service.Get${CAP_DOMAIN}ByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("${CAP_DOMAIN} not found", "error", err)
		appErr := sharedErr.NewNotFound("${CAP_DOMAIN} not found")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "${CAP_DOMAIN} fetched successfully", "data": entity})
}

func (h *Handler) Update${CAP_DOMAIN}(c *gin.Context) {
	id := c.Param("id")
	var req Update${CAP_DOMAIN}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid request body")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid request body")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	entity, err := h.service.Update${CAP_DOMAIN}(c.Request.Context(), id, req)
	if err != nil {
		h.logger.Error("Failed to update $DOMAIN", "error", err)
		appErr := sharedErr.NewDatabaseError("Failed to update $DOMAIN")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "${CAP_DOMAIN} updated successfully", "data": entity})
}

func (h *Handler) Delete${CAP_DOMAIN}(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete${CAP_DOMAIN}(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete $DOMAIN", "error", err)
		appErr := sharedErr.NewDatabaseError("Failed to delete $DOMAIN")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "${CAP_DOMAIN} deleted successfully", "data": nil})
}

func (h *Handler) List${CAP_DOMAIN}s(c *gin.Context) {
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "10")

	skip, err := strconv.ParseInt(skipStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid skip parameter", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid skip parameter")
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid limit parameter", "error", err)
		appErr := sharedErr.NewBadRequest("Invalid limit parameter")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	entities, err := h.service.List${CAP_DOMAIN}s(c.Request.Context(), skip, limit)
	if err != nil {
		h.logger.Error("Failed to list ${DOMAIN}s", "error", err)
		appErr := sharedErr.NewDatabaseError("Failed to list ${DOMAIN}s")
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "${CAP_DOMAIN}s fetched successfully", "data": entities})
}
EOF

# 7. routes.go
cat > internal/$DOMAIN/${DOMAIN}Routes.go <<EOF
package $DOMAIN

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(router *gin.Engine) {
    group := router.Group("/${DOMAIN}s")
    {
        group.POST("", h.Create${CAP_DOMAIN})
        group.GET("", h.List${CAP_DOMAIN}s)
        group.GET(":id", h.Get${CAP_DOMAIN})
        group.PUT(":id", h.Update${CAP_DOMAIN})
        group.DELETE(":id", h.Delete${CAP_DOMAIN})
    }
}
EOF

# 8. Insert into di/container.go
# Import
# sed -i '' "/import (/a\\
# 	\"project-structure/internal/$DOMAIN\"\
# " di/container.go

# # Repository
# sed -i '' "/userRepo :=/a\\
# 	${DOMAIN}Repo := $DOMAIN.NewMongoRepository(mongoDB)\
# " di/container.go

# # Service
# sed -i '' "/userService :=/a\\
# 	${DOMAIN}Service := $DOMAIN.NewService(${DOMAIN}Repo)\
# " di/container.go

# # Handler
# sed -i '' "/userHandler :=/a\\
# 	${DOMAIN}Handler := $DOMAIN.NewHandler(validate, log, ${DOMAIN}Service)\
# " di/container.go

# # Register routes
# sed -i '' "/userHandler.RegisterRoutes(router)/a\\
# 	${DOMAIN}Handler.RegisterRoutes(router)\
# " di/container.go
FILE="di/container.go"

awk -v domain="$DOMAIN" -v Domain="$PASCAL_CASE_DOMAIN" '
BEGIN {
  importInserted = 0
  repoInserted = 0
  serviceInserted = 0
  handlerInserted = 0
  routeInserted = 0
}

{
  print

  # Import line
  if (!importInserted && $0 ~ /^import \($/) {
    print "\t\"project-structure/internal/" domain "\""
    importInserted = 1
  }

  # Repository line
  if (!repoInserted && $0 ~ /\.New.*Repository\(/) {
    print "\t" domain "Repo := " domain ".New" Domain "Repository(mongoDB)"
    repoInserted = 1
  }

  # Service line
  if (!serviceInserted && $0 ~ /\.New.*Service\(/) {
    print "\t" domain "Service := " domain ".New" Domain "Service(" domain "Repo)"
    serviceInserted = 1
  }

  # Handler line
  if (!handlerInserted && $0 ~ /\.New.*Handler\(/) {
    print "\t" domain "Handler := " domain ".New" Domain "Handler(validate, log, " domain "Service)"
    handlerInserted = 1
  }

  # Route registration
  if (!routeInserted && $0 ~ /\.RegisterRoutes\(router\)/) {
    print "\t" domain "Handler.RegisterRoutes(router)"
    routeInserted = 1
  }
}

' "$FILE" > "$FILE.tmp" && mv "$FILE.tmp" "$FILE"


echo "Scaffolded CRUD and DI for domain: $DOMAIN" 