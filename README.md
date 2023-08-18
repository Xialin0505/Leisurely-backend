# Leisurely-backend

### Introduction
This is the backend of Leisurely. Leisurely is an itinerary day-trip planner app 
that generate an at most one day trip plan that matches user preference.

This backend system can be deployed on multiple server connected via a load balancer to 
achieve the best performance. It contains docker image for Postgresql database, middleware
written in Go and NLP module written in Python

### App sample
The home page, user page are shown below
<p align="center">
  <img src="./sample/main.png" width="250" title="home page">
  <img src="./sample/user profile.png" width="250" title="user profile">
</p>

The trip can be generated using the following page

<p align="center">
  <img src="./sample/new trip1.png" width="250" title="trip generation">
  <img src="./sample/new trip 2.png" width="250" title="trip generation">
</p>

User can choose the tags attached to the plan, budget, location anywhere in the world, by public transit or by driving
to generate the trip plan, or choose minimum input (date, location and the time interval they want to go out).


After the plan is generated, user can rearrange each place they are going, and also choose other place to go from
a list of alternatives that matches the user's preference.

<p align="center">
  <img src="./sample/edit trip.png" width="200" title="edit trip">
  <img src="./sample/Event detail.png" width="200" title="trip detail">
  <img src="./sample/time picker.png" width="200" title="rearrage trip">
  <img src="./sample/Alternative.png" width="200" title="alternatives trip">
</p>

Once the trip plan is confirmed, the trip will be saved to home page with a theme image attached to it.
