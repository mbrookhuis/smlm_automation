# smlm_automation
Automation for SMLM 5.x to manage all components 



Design

Data should be in a yaml file. This yaml file will contain the same sections a below.

Parameter: 
* -f <config.yaml>
* --file <config.yaml>

Following options:
* general:
  * set loglevel and logdirectory
* system:
  * organization 
  * scc credentials
  * admin user and password
* users
  * user + password or PAM
  * role
* channels
  * scc channels
  * repositories
  * channels (e.g. custom-repo, channels for repositories)
  * content life cycle (including filter)
  * channel administrators
* systemgroups
  * group
  * systems assigned
  * group administrators
  * configuration channels
  * formulars --> will follow later
* activationkeys
  * key
  * system group
  * configuration channels
* configuration channels --> will follow later
  * channel --> only salt
  * content 
  * configurationchannel administration
* images
  * images
  * creation of images (see what we did for marting)
  * assign build server
  * image administrator
* all
  * all of the above

Implementation:
* use cobra
* use viber for configuration

Example config.yaml
```yaml
general: 
  log:
    screen_level: info
    file_level: debug
    file_path: /var/log/do_smlm/do_smlm.log

system:
  organization_name: orgname
  scc:
    user_1: password_1
    user_2: password_2
  admin_user: sm-admin
  admin_password: SUSE4ever!



```