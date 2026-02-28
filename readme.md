# Envoy | PostgreSQL Management Tool

## Why did I decide to create this project?

I've worked for startups and enterprises, and I've seen the same problems over and over:

- You need to perform a database migration between environments and you have to go over 5 different tools and terminals to do it
- You need to manage multiple databases and you don't know which one is which
- You need to double-check permissions in PgAdmin and psql to make sure you have the right permissions in your schema
- You need to deal with TLS certificates and connection strings
- The version control and audit log is poor or nonexistent when dealing with database changes

### Example situation

1. It's Friday at 10am, you need to deploy the changes from staging to production
2. You get the production database credentials  
3. You replace an environment variable in your .env file
4. You're running the migration command
5. You verify the migration was successful
6. You try the changes in production
7. You forgot to give permissions to the new tables
8. You have to go back to PgAdmin or psql to give permissions
9. You try again in the web application and it works
10. Oh congrats, you have to same thing next week across all environments if you work for a startup!

## Who's this tool for?

- Startup engineers who deal with multiple databases and need a unified interface
- Teams that want to standardize their database management processes
- Consultants who need to quickly understand and manage multiple database environments
- Solo developers managing multiple database environments

## Notes

This project is not trying to replace tools like PgAdmin, DBeaver, or other database management tools. It's trying to provide a unified interface for managing multiple database environments and performing database migrations across them.

I'm giving the best support to PostgreSQL databases, but it can be extended to other databases in the future.

## Clarification and pricing

This is a complete free and open source project. My goal is to help developers and small teams manage their database environments without the pain of dealing with multiple tools and terminals. I have created scripts to deal with the complexity of database management, but I wanted to make database management easier for everyone.

## How is this different from other tools?

Envoy is opinionated and made for the startup world. I'm not competing against Atlas, Liquibase, or other database migration tools. They're corporate and great tools to manage database migrations. This is for small team performing DevOps in multiple environments without the management overhead, you have everything in a unified interface. Bash scripts and env variables are too messy, Atlas and Liquibase are too complex for startups that need to deploy changes quickly.