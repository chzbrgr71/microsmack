# Microsmack Demo App

This is a simple set of Go apps to demonstrate web and microservices API communication. These can be used as the basis for some simple container orchestration demos (Eg - Kubernetes)

## Applications

### Web UI
  - Calls web API using ENVVAR
  - Listens on port 8080  

  ```docker build --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` --build-arg VCS_REF=`git rev-parse --short HEAD` --build-arg VERSION=$VERSION -t chzbrgr71/smackweb .```

  ```helm install --name=smackweb ./charts/smackweb```
  
#### Container Details
[![](https://images.microbadger.com/badges/image/chzbrgr71/smackweb.svg)](https://microbadger.com/images/chzbrgr71/smackweb "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/chzbrgr71/smackweb.svg)](https://microbadger.com/images/chzbrgr71/smackweb "Get your own version badge on microbadger.com")
[![](https://images.microbadger.com/badges/commit/chzbrgr71/smackweb.svg)](https://microbadger.com/images/chzbrgr71/smackweb "Get your own commit badge on microbadger.com")

### Web API
  - Listens on port 8081

  ```docker build --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` --build-arg VCS_REF=`git rev-parse --short HEAD` --build-arg VERSION=$VERSION -t chzbrgr71/smackapi .```
  
  ```helm install --name=smackapi ./charts/smackapi```

#### Container Details
[![](https://images.microbadger.com/badges/image/chzbrgr71/smackapi.svg)](https://microbadger.com/images/chzbrgr71/smackapi "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/chzbrgr71/smackapi.svg)](https://microbadger.com/images/chzbrgr71/smackapi "Get your own version badge on microbadger.com")
[![](https://images.microbadger.com/badges/commit/chzbrgr71/smackapi.svg)](https://microbadger.com/images/chzbrgr71/smackapi "Get your own commit badge on microbadger.com")