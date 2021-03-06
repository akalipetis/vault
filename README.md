VAULT  [![Build Status](https://travis-ci.org/franela/vault.svg?branch=master)](https://travis-ci.org/franela/vault)
=====


# What is Vault?

Vault is a __secure__ cross-platform KV store wrapped around GPG which makes team collaboraion extremely easy. 
You can define as many Vaults as needed and use them rightaway without changing any configuration at all.


# Why Vault?

1. Because security is something really important for us and we take it seriously
2. It really simplifies the friction to work with GPG between multiple teams and recipients
3. We needed something which was easy to use and platform agnostic
4. It was fun to do it and to learn during the process


# What can I do with Vault?

Using Vault is super easy. Below is the list of all the possible operations:

`vault init` - Creates a new Vault in the directory you're located  
`vault add` - Adds one or more recipients to your Vault. Vault will automatically re-encrypt all your files for the new recipient  
`vault set` - Stores something (text or file) into your Vault and encrypts it for all of your vault recipients  
`vault get` - Retrieves an encrypted file from your Vault and decrypts it.  
`vault edit` - Decrypts a file and opens it in the default editor and re-encrypts it after save and exit.
`vault remove` - Removes specified recipients from the Vault. It will automatally keep te integrity of your Vault by upgrading encrypted recipients  
`vault import` - Imports all the current recipients in your vault. Alternatively you can specifiy a keyserver to import from  
`vault recipients` - Lists vault recipients and it's current status in the user keyring (missing / untrusted / trusted) i.e:  

![recipients](https://cloud.githubusercontent.com/assets/1578458/10444397/5e543a6a-713c-11e5-98da-99dfd38e3b88.png)


# This sounds cool, where do I get it?


You can find Vault for your favorite distrution below:  


<p align="center">
  <a href="https://github.com/franela/vault/releases/download/0.0.2/linux.zip" ><img width="150px" height="150px" src="http://imagenes.es.sftcdn.net/blog/es/2013/09/Tux-Seguridad.png" alt="Linux"/> </a>
  <a href="https://github.com/franela/vault/releases/download/0.0.2/windows.zip" ><img width="150px" height="150px" src="http://webpamplona.com/wp-content/uploads/2014/06/windows.png" alt="Windows" /> </a>
  <a href="https://github.com/franela/vault/releases/download/0.0.2/darwin.zip" ><img width="150px" height="150px" src="http://www.weblinear.fr/img/glyphicons/macos-logo.jpg" alt="MacOS" /> </a>
  <a href="https://github.com/franela/vault/releases/download/0.0.2/freebsd.zip" ><img width="130px" height="130px" src="http://1.bp.blogspot.com/-mls96EYcCoA/U-sS1D6FknI/AAAAAAAATqk/BCRJYO9jR4U/s1600/freebsd.png" alt="BSD" /> </a>
  
</p>


More distributions and source code can be found [here](https://github.com/franela/vault/releases)

## Installing

Download the correct binary for your distribution (clck images above) into your PATH and check that [GPG](https://www.gnupg.org/download/) installed.


# Contributing

We're always working to improve vault, but if you find a bug or you just want to collaborate you can send us your PR (include tests when possible).

