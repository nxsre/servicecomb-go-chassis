package servicecomb

// Basic API
var HEALTH = "/v4/default/registry/health"
var VERSION = "/v4/default/registry/version"

// Micro-Service API's
var CHECKEXISTENCE = "/v4/default/registry/existence"
var GETALLSERVICE = "/v4/default/registry/microservices"
var GETSERVICEBYID = "/v4/default/registry/microservices/:serviceId"
var REGISTERMICROSERVICE = "/v4/default/registry/microservices"
var UPDATEMICROSERVICE = "/v4/default/registry/microservices/:serviceId/properties"
var UNREGISTERMICROSERVICE = "/v4/default/registry/microservices/:serviceId?force=1"
var GETSCHEMABYID = "/v4/default/registry/microservices/:serviceId/schemas/:schemaId"
var UPDATESCHEMA = "/v4/default/registry/microservices/:serviceId/schemas/:schemaId"
var UPDATESCHEMAS = "/v4/default/registry/microservices/:serviceId/schemas"
var DELETESCHEMA = "/v4/default/registry/microservices/:serviceId/schemas/:schemaId"
var CREATEDEPENDENCIES = "/v4/default/registry/dependencies"
var GETCONPRODEPENDENCY = "/v4/default/registry/microservices/:consumerId/providers"
var GETPROCONDEPENDENCY = "/v4/default/registry/microservices/:providerId/consumers"

// Instance API's
var FINDINSTANCE = "/v4/default/registry/instances"
var GETINSTANCE = "/v4/default/registry/microservices/:serviceId/instances"
var GETINSTANCEBYINSTANCEID = "/v4/default/registry/microservices/:serviceId/instances/:instanceId"
var REGISTERINSTANCE = "/v4/default/registry/microservices/:serviceId/instances"
var UNREGISTERINSTANCE = "/v4/default/registry/microservices/:serviceId/instances/:instanceId"
var UPDATEINSTANCEMETADATA = "/v4/default/registry/microservices/:serviceId/instances/:instanceId/properties"
var UPDATEINSTANCESTATUS = "/v4/default/registry/microservices/:serviceId/instances/:instanceId/status"
var INSTANCEHEARTBEAT = "/v4/default/registry/microservices/:serviceId/instances/:instanceId/heartbeat"
var INSTANCEWATCHER = "/v4/default/registry/microservices/:serviceId/watcher"
var INSTANCELISTWATCHER = "/v4/default/registry/microservices/:serviceId/listwatcher"

//Governance API's
var GETGOVERNANCESERVICEDETAILS = "/v4/default/govern/microservices/:serviceId"
var GETRELATIONGRAPH = "/v4/default/govern/relations"
var GETALLSERVICEGOVERNANCEINFO = "/v4/default/govern/microservices"

//Rules API's
var ADDRULE = "/v4/default/registry/microservices/:serviceId/rules"
var GETRULES = "/v4/default/registry/microservices/:serviceId/rules"
var UPDATERULES = "/v4/default/registry/microservices/:serviceId/rules/:rule_id"
var DELETERULES = "/v4/default/registry/microservices/:serviceId/rules/:rule_id"

//Tag API's
var ADDTAGE = "/v4/default/registry/microservices/:serviceId/tags"
var UPDATETAG = "/v4/default/registry/microservices/:serviceId/tags/:key"
var GETTAGS = "/v4/default/registry/microservices/:serviceId/tags"
var DELETETAG = "/v4/default/registry/microservices/:serviceId/tags/:key"

// HTTP METHODS
var GET = "GET"
var POST = "POST"
var UPDATE = "PUT"
var DELETE = "DELETE"
