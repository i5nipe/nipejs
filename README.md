# JSScan

> Ler lista de arquivos JS e procura dados sensiveis via regex. 

### Ajustes e melhorias

O projeto ainda está em desenvolvimento e as próximas atualizações serão voltadas nas seguintes tarefas:

- [x] Ler lista de URLs via STDIN
- [ ] Modificar o tempo de timeout
- [ ] Concorrência para bater as regex

## ☕ Usando JSScan

Para usar JSScan, siga estas etapas:

```
jsscan -urls jsfile

jsscan -urls ~/Path/to/jsfile -s 

cat jsfile | jsscan
```



