A proposta deste repositório é mostrar um exemplo de API feita em GO usando algumas abordagens que eu considero essenciais:

#### Pipeline de CI com Github Actions
Antes de fazer o merge de branch é importante verificar a qualidade do código atraves de ferramenas que analisam aspectos de segurança, vulnerabilidades de bibliotecas e padrões de código.
Também é importante rodar os testes de unidade, integração e E2E. 

#### Devcontainer
Ambiente de desenvolvimento isolado e padronizado, garantindo que todos os desenvolvedores tenham as mesmas dependências e ferramentas configuradas corretamente.

#### Entidades de Domínio
Realocar regras de negócio que estão espalhadas em diversos serviços, e centralizá-las em entidades de domínio.
Basicamente utilizei uma caracteristica do DDD sem necessariamente usar todo o poder do DDD.

#### Testes de Integração com Testcontainers
Com testcontainers é possível criar containers dinamicamente para fazer testes em implementações concretas de serviços externos.
Exemplo: nessa aplicação eu faço a utilização do aws secrets manager. Com o testcontainer eu consigo fazer teste de integração nessa implementação.

#### Testes de E2E
Uma pratica interessante é subir uma infraestrutura com Terraform o mais próximo possível da PROD e rodar os testes E2E nessa infra. Depois, basta destruir essa infraestrutura automaticamente no final do processo.

#### Testes de Unidade
#### Logs
