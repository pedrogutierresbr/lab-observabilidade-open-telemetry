## Lab Observabilidade e Open Telemetry

<br>

## Sobre o projeto

Este é o repositório destinado ao laboratório de Observabilidade e Open Telemetry do curso Pós Goexpert da faculdade FullCycle. O projeto permite ao usuário visualizar a temperatura atual e a cidade, e também monitorar a aplicação usan Otel + Zipkin.

<br>

## Funcionalidades

-   Receber um CEP via requisição POST.
-   Consultar a API ViaCEP para identificação do cep indicado;
-   Consultar a temperatura na localização indicada, pela API WeatherAPI;
-   Visualizar a temperatura em diversas unidades de medida, junto com a cidade.

<br>

## Como executar o projeto

### Pré-requisitos

Antes de começar, você vai precisar ter instalado em sua máquina as seguintes ferramentas:

-   [Git](https://git-scm.com)

-   [VSCode](https://code.visualstudio.com/)

-   [Postman](https://www.postman.com/)

-   [Docker](https://www.docker.com/)

-   [WEATHER_API_KEY](https://www.weatherapi.com/) (Necessário cadastro para gerar key)

<br>

#### Acessando o repositório

```bash

# Clone este repositório

$ git clone https://github.com/pedrogutierresbr/lab-observabilidade-open-telemetry.git

```

<br>


#### Executando a aplicação em ambiente dev

```bash

# Crie um arquivo .env na raiz do projeto.

# Adicione as seguintes chaves no arquivo .env:

	URL_ZIPKIN=localhost
	URL_SERVICE_B=http://localhost:8081/cep/
	URL_VIACEP=http://viacep.com.br/ws/
	URL_WEATHERAPI=http://api.weatherapi.com/v1/current.json
	WHEATER_API_KEY={WheaterAPI key}

# Abrir um terminal

# Executar a aplicação

$ docker-compose up -d

# Para pausar a aplicação

$ docker-compose down

```

<br>

#### Realizando requisição localmente

Serviços estarão disponíveis na seguinte porta:

-   Web Server : 8080

##### Web Server

```bash

# Crie uma requisição com o auxilio do Postman

POST http://localhost:8080/cep

body
{
	"cep": "20021010"
}

```

<br>

#### Rastreamento

Esta aplicação utiliza OpenTelemetry e Zipkin:

```bash

$ http://localhost:9411

```

<br>

## Licença

Este projeto esta sobe a licença [MIT](./LICENSE).

Feito por Pedro Gutierres [Entre em contato!](https://www.linkedin.com/in/pedrogabrielgutierres/)
