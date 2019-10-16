# PRR - Laboratoire 1 : Synchronisation d'horloge

## Utilisation

### Sur la machine hôte

### Avec Docker

Afin de lancer le maître, il suffit d'entrer la commande suivante.

```
docker-compose up -d clock-master
```

Cela téléchargera le container `golang` avant de copier les fichiers sources dans le container, le construire et le lancer en arrière-plan.

Ensuite, pour lancer les esclaves et afficher les logs, entrer cette commande :

```
docker-compose up --scale clock-slave=5
```

Si vous préférez lancer les containers en arrière plan, ajouter le flag `-d` et utilisez ensuite la commande

```
docker-compose logs
```

pour voir la sortie des programmes.

## Fonctionnement

## Problèmes connus

- Dans le slave, pas de timeout lors de l'attente de DELAY_RESPONSE.

- Nous n'avons pas trouvé la source d'une erreur arrivant de manière peu reproductible "unexpected EOF".
