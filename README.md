# secretSanta: simple app for playing Secret Santa
* Do you want to play the Secret Santa with your friends, coworkers, family?
* You or your friends have a wife or a husband, a girlfriend or a boyfriend 
who is going to play?
* Do you want to restrict their gift recipients to keep an interest?
  After all, your friends or coworkers will have gone to please their soulmate,
  and it doesn't depend on playing the game.
* Online services, chatbots and other applications don't let to set up such
  restrictions?

Probably, this application will help you!
It allows setting a restriction list for each player
to exclude matching written participants if you want.

To imagine this app as external human who are asked to random match
gift recipients for each player and to send email with results them.

Requirements
--
* Each participant must have a valid email address because the result is
needing to send him directly.
* One participant is one email address.
* You must have active email address, which supports sending via SMTP.

Limits
--
* By design and implementation, application can send results only by email.
* Email text templates support only russian language currently.

How to use?
--
The program is the console application.
Download an archive with your OS version, extract ready app. 

To generate example config file run:

    ./secretSanta createConfig

File `conf.toml` will appear.
Open it in any text redactor and set up configuration.

Then try to run application in test mode:

    ./secretSanta dryRun

It runs matching recipients for players, does imitation sending emails
and prints result in the console.
It is useful to run before direct sending to check,
was the matching correct with entered restrictions.

To match and actually to send a result run:

    ./secretSanta run

After sending, the result won't be printed in the console.

What about config file? What it has and how to make settings?
--
Config file has `SmtpHost` and `SmtpPort` values.
It is a host and port to send email via SMTP.
For example, there is yandex host with port 587
in the [conf-example.toml](conf-example.toml).
You can use other values.

Then you should set a valid sender data: 
`EmailSender` - email address which is used to send email,
`NameSender` - sender name, which views for each recipient,
`PasswordSender` - sender email password.
Usually it's using individual mail applications with specific password,
e.g. [Yandex](https://yandex.com/support/mail/mail-clients/others.html#imap__imap-app-pass

So, you should make mapping email participants with their names.
Names are used in an email body. 
Email recipient will see the player's name, whom he will prepare a gift.

It remains to set restrictions. 
You should make mapping participant emails with an array of emails
other players, which must never be matched to the current player.
For example, this is saying:

    "email1@test.com" = [ "email2@test.com" ]
it is guaranteed player with `email2@test.com` email never
get caught to player with `email1@test.com` email address.

You can list restrictions using comma, e.g.:

    "email6@test.com" = [ "email5@test.com", "email1@test.com" ]
If limits get a situation, that it is impossible
to find a recipient at least for a one player, the application will exit with error
and print the player, who has no recipients.

If you don't want to set limits, you can set an empty array, e.g.:

    "email7@test.com" = []      
or entirely leave a blank restriction section.

Enjoy! 😏
