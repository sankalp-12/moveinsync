# MoveInSync Case Study: Smart Cab Allocation System

This repository contains the case study analysis report, the demo video as well as the code implementation for the MoveInSync case study titled "Smart Cab Allocation System". 

Click [here](https://github.com/sankalp-12/moveinsync/blob/master/Analysis%20Report.pdf) for the case study analysis report. <br/>
Click [here](https://drive.google.com/file/d/1JqV-8VlDCL_KdXBnBOIRIT6FE7_JAv_g/view?usp=sharing) for the demonstration video of the implementation. <br/>
Click [here](https://github.com/sankalp-12/moveinsync) for the code repository of the implementation.

## Tech-Stack

- **Languages:** _Go_; for backend development & _YAML_; for writing config files.
- **Databases:** _MongoDB_: used _[MongoDB Atlas](https://www.mongodb.com/atlas/database)_ to manage cloud instances.
- **Deployment Tools:** _Docker_; used _[Docker-Compose](https://docs.docker.com/get-started/08_using_compose/)_ to manage & deploy multi-container applications. [Postman](https://www.postman.com/product/what-is-postman/); for API-testing. [Prometheus](https://prometheus.io/) & [Grafana](https://grafana.com/); as the monitoring stack.

## Repository Structure

- `admin-service`: Service responsible for all admin actions; namely: admin creation, admin login, adding cabs & suggesting/allocating cabs.
- `cab-data-service`: Service responsible for real-time cab location data integration. It is a _web-socket_ server that listens for cab devices reporting location data and writes it back to the database.
- `user-service`: Service responsible for all user actions; namely: user creation, user login, booking trips & requesting to display engaged cabs.
- `prometheus`: Responsible for metrics collection from the necessary services.
- `grafana`: Responsible for metrics visualizations as well as alerting for the necessary services.

## How to run locally

- Clone the GitHub repository using the following command to your local machine.
	```
	git clone https://github.com/sankalp-12/moveinsync
 	```

- Make sure you have installed _Docker_ and _Docker-Compose_ on your machine. Now, change the working directory to the newly cloned repository.
	```
	cd moveinsync
 	```

- Now, to start all the containers, run the _docker-compose.yml_ file.
   	```
    sudo docker-compose build && sudo docker-compose up
   	```
   If you face any errors during the executions of the above command, ensure that all the ports required are free on your machine. You can check for the required ports in the `docker-compose.yml` file.

- After execution, run the command: `docker ps` to ensure that all the containers are up and running.
   
- Once verified, access the port: 3001 by _default_ assigned to the **Grafana** instance. The credentials by default are: 
```
Username: admin
Password: admin
```

- Once logged in, you can start making dashboards to monitor the application metrics and send alerts. A template dashboard has been provided in `grafana/dashboards/moveinsync.json` with a limited number of visualizations, which you can upload to **Grafana** through the _import_ option while creating a new dashboard.

- The setup is now complete! You can now start sending HTTP requests to the services (using _curl_ or _Postman_). The different endpoints and their respective request body structures are provided below.

## API Documentation

### Admin-Service

- _Create Admin_
   ```yaml
    - Endpoint: `http://localhost:8081/api/v1/admin/create`
    - Method: POST
    - Request Body:
      {
          "username": "demo",
          "password": "moveinsync"
      }
    - Expected Response:
      {
          "status": "success"
      }
   ```

- _Admin Login_
   ```yaml
    - Endpoint: `http://localhost:8081/api/v1/admin/login`
    - Method: POST
    - Request Body:
      {
          "username": "demo",
          "password": "moveinsync"
      }
    - Expected Response:
      {
          "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NzUxMzcsInVzZXJuYW1lIjoic2Fua2FscCJ9.0ySKDMXCGP7mFytkrAFFQo2JonX955OlKlWClwbTHLw"  
      }
   ```

 - _Add Cabs_
   ```yaml
    - Endpoint: `http://localhost:8081/api/v1/admin/addcabs`
    - Method: POST
    - Authorisation Header: Bearer [token]
    - Request Body:
      {
          "location":
          {
              "type": "Point",
              "coordiantes": ["[longitude]", "[latitude]"]
          }
          "status": "Available"/"Busy"
      }
    - Expected Response:
      {
          "status": "success"  
      }
   ```

### User-Service

- _Create User_
   ```yaml
    - Endpoint: `http://localhost:8080/api/v1/user/create`
    - Method: POST
    - Request Body:
      {
          "username": "demo",
          "password": "moveinsync"
      }
    - Expected Response:
      {
          "status": "success"
      }
   ```

- _User Login_
   ```yaml
    - Endpoint: `http://localhost:8080/api/v1/user/login`
    - Method: POST
    - Request Body:
      {
          "username": "demo",
          "password": "moveinsync"
      }
    - Expected Response:
      {
          "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NzUxMzcsInVzZXJuYW1lIjoic2Fua2FscCJ9.0ySKDMXCGP7mFytkrAFFQo2JonX955OlKlWClwbTHLw"  
      }
   ```

 - _Book Trip (Available Cabs)_
   ```yaml
    - Endpoint: `http://localhost:8080/api/v1/user/booktrip`
    - Method: POST
    - Authorisation Header: Bearer [token]
    - Request Body:
      {
          {
              "longitude": "81.60",
              "latitude": "21.24"
          }
      }
    - Expected Response:
      {
          {
              "_id":"661f75e1b80a04c2f7adc774",
              "distance":29187.820383017817,
              "last_updated":"2024-04-17T07:10:25.83Z",
              "location":
              {
                  "coordinates": [81.31886422688906,21.249442177652497],
                  "type":"Point"
              },
              "status":"Available"
          }
      }
   ```
   
- _Display Nearby Cabs (Busy Cabs)_
   ```yaml
    - Endpoint: `http://localhost:8080/api/v1/user/displaynearbycabs`
    - Method: POST
    - Authorisation Header: Bearer [token]
    - Request Body:
      {
          {
              "longitude": "81.60",
              "latitude": "21.24"
          }
      }
    - Expected Response:
      {
          {
              "_id":"661f75e1b80a04c2f7adc774",
              "distance":29187.820383017817,
              "last_updated":"2024-04-17T07:10:25.83Z",
              "location":
              {
                  "coordinates": [81.31886422688906,21.249442177652497],
                  "type":"Point"
              },
              "status":"Busy"
          }, ... (upto 5 cabs)
      }
   ```

### Cab-Data Service

- _Real-Time Cab-Location Data Integration_
   ```yaml
    - Endpoint: `http://localhost:8082/cab/ws`
    - Method: POST
    - Request Body:
      { 
          "ID": "cab_id",
          "location":
          {
              "type": "Point",
              "coordiantes": ["[longitude]", "[latitude]"]
          }
          "status": "Available"/"Busy"
      }
   ```