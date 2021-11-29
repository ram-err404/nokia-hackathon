# Nokia-Hackathon
We developed an auto scaling micro-service that scale up or down the pods based on the increase/decrease in API requests from mediation xApps.

### Team Name: 
> The A-Team
### Team Members: 
> Supriya Pal, Srinidhi Chaluvaraju, Prabhu Swamy NM and Ramji Misra.


## Background of our Idea:
> Nokia proprietary and value-added services in the SEP system are mediation xApps. Mediation xApps exposes set of APIs towards use case xApps which can be Nokia or 3rd party developed which acts as clients. 

> Number of clients could grow over time, and it could proportionally increase the load and expects better reliability on the mediation xApps

#### What is the value-add for SEP platform?
> Customized monitoring of mediation xApp APIs usage (APIs/sec) in the SEP common dashboard. This could be eventually linked to the licensing

>  Need based deployment of mediation xApps

>  Horizontal Autoscaling based on API load 

>  Mediation xApp API usage vs RAN/UE behaviours can be learnt, and it could be  used  for Self Organized Platform with the aid of AI/ML framework


## Idea Concept:
>1. API Gateway is the single endpoint for all client services hence all API transaction goes via API Gateway. This is the source of data for the situation analysis for our idea.

>2. Ingenious Service is the brain of our idea. It performs following activities:

>>a) Periodically (1 sec) fetches API Gateway metrics and parse the data to extract API count per mediation xApp.

>>b) Create customer metrics  (API rate / sec) and send to Prometheus, which can be viewed in Grafana 

>>c) If API rate crosses the configured threshold value. It scale out corresponding mediation xApp and scale in if API rate falls lower than 25% of threshold values.


#### For more detailed description refer the attached ppt present at: 
