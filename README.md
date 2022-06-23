# service_catalog

## Description
This is an example of implementation of Service catalog


## Description

Following the Clean Architecture by Uncle Bob
 * Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
 * Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
 * Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
 * Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
 * Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

This project has  4 Domain layer :
 * Domain Layer
 * Repository(DB) Layer
 * Usecase(Service) Layer
 * Handler(Controller) Layer


Other assumptions
  * A default version is created whenever a service is created
  * Creating another version of a service calls the create version endpoint
  * The direction of sort should be in descending order
  * Items on each page limited to 12
  * Used entgo orm for easy database access
