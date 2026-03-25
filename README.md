# 🧰 Chest

A minimalist command-line password manager to store secrets and manage remote connections in one place.

### 🔐 Secret Management
Store and retrieve different kinds of sensitive data:
- Usernames & Passwords
- API Tokens
- Environment variables

### 🚀 Remote Access
Quickly connect to your infrastructure:
- **SSH/Login:** Custom login commands and scripts.
- **Kubernetes:** Manage and switch between multiple `kubeconfig` files seamlessly.

---
*Status: Work in Progress 🛠️*

###### USAGE

**chest create**

invoca la richiesta di tipo tra i tipi supportati(chiavi della mappa dei chest parser)



, nome, tipo e descrizione definita nel new del base chest
invoca la richiesta del jewel di decript definita nel chest di quel tipo e crea il chest con quel tipo

**chest open <name>** to open a chest

verifica che il chest ci sia
legge il chest type dal json
chiede di inserire il jewel secondo il metodo di autenticazione previsto dal jewel
prova a aprire il json lanciando la goruotine di decrypt di quel metodo
create o update a file in /run/user/<userid> criptato con la chiave machine-id/userid/shellid in cui è salvato il json
[{
    chest_name: <nome del chest>
    chest_type: <tipo del chest>
    description: <descrizione del chest>
    jewels: [{
        <jewel to open the chest>
    }]

}]

**chest ls** to list all chest
lancia una goroutine per ogni chest con la funzione printChest

**chest close** to close the chest
**chest add <name>** to add a chest
**chest delete <name>** to delete a chest


#### esample for kube usage

**chest kube ls** to list all kube jewel
**chest kube add <name>** to add a kube jewel
**chest kube delete<name>** to remove a kube jewel
**chest kube edit <name>** to edit a kube jewel
**chest kube copy <name>** to copy the kube jewel
**chest kube print <name>** to print the kube jewel
**chest kube <name>** to use a jewel