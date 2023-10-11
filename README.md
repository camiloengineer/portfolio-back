# Backend de Portafolio

Este proyecto contiene el backend de mi [portafolio profesional](https://www.camiloengineer.com/), implementado utilizando Go. 

## Tecnologías Utilizadas

- **Lenguaje**: Go
- **Frameworks y ORM**: Gorilla/mux, Gorm ORM, Air
- **Base de Datos**: PostgreSQL
- **Cloud**: GCP SDK, Cloud SQL
- **Otras Tecnologías**: Docker

## Cómo Ejecutar Localmente

### Prerequisitos

- Go ([descarga e instalación](https://golang.org/dl/))

### Instrucciones

1. **Clona el Repositorio**

   ```bash
   git clone https://github.com/camiloengineer/portfolio-back.git
   ```

2. **Ejecuta la Aplicación**

   ```bash
   go run .
   ```

### Modelo de Datos


1. **Project**

| Column          | Type     | Description                            |
|-----------------|----------|----------------------------------------|
| id              | SERIAL   | Unique identifier for the project      |
| url             | TEXT     | URL related to the project             |
| image           | TEXT     | Reference to the project image         |
| is_professional | TEXT     | Reference to the project image         |


2. **ProjectTranslations**

| Column       | Type     | Description                             |
|--------------|----------|-----------------------------------------|
| id           | SERIAL   | Unique identifier for the translation   |
| project_id   | INT      | Foreign key to `Projects` table         |
| language     | TEXT     | Language code (e.g., 'en', 'es')        |
| title        | TEXT     | Translated title of the project         |
| description  | TEXT     | Translated description of the project   |

3. **Categories**

| Column       | Type     | Description                            |
|--------------|----------|----------------------------------------|
| id           | SERIAL   | Unique identifier for the category     |
| name         | TEXT     | Name of the category                   |

4. **ProjectCategories**

| Column       | Type     | Description                            |
|--------------|----------|----------------------------------------|
| project_id   | INT      | Foreign key to `Projects` table        |
| category_id  | INT      | Foreign key to `Categories` table      |


### Example Data


1. **Project**

| id | url                  | image    | is_professional |
|----|----------------------|----------|-----------------|
| 1  | https://example.com  | image1   | false           |
| 2  | https://example2.com | image2   | true            |



2. **ProjectTranslations**

| id | project_id | language | title    | description                 |
|----|------------|----------|----------|-----------------------------|
| 1  | 1          | en       | Example  | This is an example project. |
| 2  | 1          | es       | Ejemplo  | Este es un proyecto ejemplo.|

3. **Categories**

| id | name       |
|----|------------|
| 1  | Web        |
| 2  | Design     |

4. **ProjectCategories**

| project_id | category_id |
|------------|-------------|
| 1          | 1           |
| 1          | 2           |
| 2          | 1           |



## ✒️ Autor

* **Camilo González** 
    * [Linkedin](https://www.linkedin.com/in/camiloengineer/)
    * [Website](https://www.camiloengineer.com/)
    * [Email](mailto:camilo@camiloengineer.com)
