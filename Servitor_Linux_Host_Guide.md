# Servitor — Guia de Preparação do Host Linux (Ubuntu Server)

Este guia prepara um host Ubuntu para laboratório e administração remota. Ubuntu é o ambiente de implantação inicial do agente Servitor; não é o nome do produto. A telemetria do agente foi projetada para sair do host por HTTPS/TLS, portanto não abra uma porta de entrada para o agente no firewall.

> Estes comandos preparam o sistema operacional. Eles não instalam nem ativam o agente Servitor.

## 1. Atualizações

```bash
sudo apt update
sudo apt upgrade -y
sudo reboot
```

Reconecte por SSH após a reinicialização.

## 2. Administração SSH

Instale e habilite o serviço:

```bash
sudo apt install openssh-server -y
sudo systemctl enable --now ssh
sudo systemctl status ssh
```

Configure autenticação por chave pública e confirme que uma segunda sessão SSH funciona antes de desativar autenticação por senha. Mantenha acesso local ou console disponível durante mudanças de rede.

## 3. Firewall UFW

Primeiro permita SSH para evitar bloquear o acesso remoto. Habilite o UFW apenas depois de validar a regra:

```bash
sudo ufw allow OpenSSH
sudo ufw status numbered
sudo ufw enable
sudo ufw status verbose
```

Não abra portas de entrada para métricas do agente. Se o host também servir uma aplicação HTTPS, abra 443/tcp somente quando essa aplicação estiver instalada e for necessária:

```bash
sudo ufw allow 443/tcp
```

## 4. Fail2Ban

```bash
sudo apt install fail2ban -y
sudo systemctl enable --now fail2ban
sudo systemctl status fail2ban
```

Use configurações locais em jail.local em vez de editar arquivos padrão do pacote. Revise logs e regras para confirmar que a política corresponde ao método de autenticação SSH adotado.

## 5. Verificação operacional

```bash
hostnamectl
ip address
sudo systemctl status ssh
sudo systemctl status fail2ban
sudo ufw status verbose
```

Registre hostname, versão do Ubuntu, disponibilidade dos sensores de temperatura (/sys/class/hwmon e /sys/class/thermal) e acesso de saída HTTPS. Sensores variam por hardware e firmware; ausência de leitura deve ser tratada como telemetria indisponível.

## Escopo do produto

A proposta do Servitor é coletar temperatura, CPU, memória e heartbeat com um agente local leve, manter buffer limitado durante falhas de rede e enviar dados ao plano de controle. Alertas térmicos, painel, workflows de incidentes, relatórios, e credenciais de agente por host estão descritos como capacidades de produto no [blueprint de arquitetura](docs/SERVITOR_PRODUCT_ARCHITECTURE_BLUEPRINT.md); este guia não afirma que essas capacidades já estejam implantadas.
