# Rules

1. Um investimento não pode existir sem estar vinculado uma conta
2. No primeiro mês de uso, o usuário apenas informa os valores iniciais dos investimentos e aportes.
3. No primeiro mês, não existe massa histórica suficiente para calcular métricas como juros médios, rendimento em dinheiro ou crescimento comparativo.
4. No primeiro mês, a visualização deve considerar apenas os valores de aporte e o total investido disponível naquele momento.
5. A partir do mês seguinte, e especialmente no primeiro dia de cada mês, o usuário deve atualizar os valores atuais de cada investimento.
6. As movimentações aceitas para evolução mensal do investimento são: `INVESTMENT_CREATED`, `CONTRIBUTION`, `INTEREST` e `ADJUSTMENT`.
7. `INVESTMENT_CREATED` representa o lançamento inicial de um investimento.
8. `CONTRIBUTION` representa adição de dinheiro feita pelo usuário no investimento.
9. `INTEREST` representa rendimento obtido pelo investimento.
10. `ADJUSTMENT` representa ajuste manual de valor para corrigir inconsistências, como erro de digitação ou ajuste excepcional informado pelo usuário.
11. O cálculo de rendimento do mês deve sempre comparar o valor atual com o valor consolidado do mês anterior.
12. O ganho em dinheiro do mês deve ser apurado com base na diferença entre o montante do mês anterior e o montante atual, desconsiderando entradas classificadas como aporte.
13. A referência para cálculo de juros nunca deve usar o próprio mês isoladamente sem considerar o valor base do mês anterior.
14. Aportes não entram no cálculo de ganho por juros.
15. Apenas movimentações classificadas como `INTEREST` devem compor o cálculo de rendimento por juros.
16. Correções não devem ser tratadas automaticamente como rendimento; elas existem para ajuste operacional e precisam permanecer separadas da apuração de juros.

## Exemplo de apuração

1. Se em julho o total investido era `10.000` e em agosto o total passou a `11.000`, existe uma variação bruta de `1.000`.
2. Essa variação não pode ser classificada integralmente como rendimento sem separar o que foi `CONTRIBUTION`, `INTEREST` e `ADJUSTMENT`.
3. Para apurar rendimento real, o backend deve considerar apenas os lançamentos do tipo `INTEREST` dentro do período.

## Impacto no dashboard

1. Quando não existir mês anterior para comparação, métricas comparativas devem vir vazias, zeradas ou com regra explícita definida pela aplicação.
2. O total de aportes do mês deve somar somente movimentações do tipo `CONTRIBUTION`.
3. O total ganho em juros do mês deve somar somente movimentações do tipo `INTEREST`.
4. O crescimento patrimonial pode considerar a variação consolidada do patrimônio no período, mas o rendimento por juros deve continuar separado dos aportes.