## Example context

```
[autodial]
exten => _XXX,1,NoOp(Local to ${PHONE} id ${OID})
    same => n,Dial(SIP/${EXTEN},60)

exten => _7XXXXXXXXXX,1,NoOp(External)
    same => n,Dial(SIP/external/${EXTEN},60)

;--
exten => h,1,GotoIF( $[ "${DIALSTATUS}" = "ANSWER" ]?end)
    same => n,GotoIF( $[ "foo${OID}" = "foo" ]?end)
    same => n,NoOp(${CURL(http://manager-url/noanswer?id=${OID}))
    same => n(end),NoOp()
--;

[autodial_answer]
exten => s,1,NoOp(Welcome to autodial ${PHONE} id ${OID})
    same => n,Wait(1)
    same => n,Playback(${SOUND})
    same => n,Hangup

;--
exten => h,1,GotoIF( $[ "foo${OID}" = "foo" ]?end)
    same => n,NoOp(${CURL(http://manager-url/answer?id=${OID}))
    same => n(end),NoOp()
--;
```

## Example request

```
curl "http://pbxhost:9002/api/originate/?oid=9fCCWE0BHSWAtUfe&phone=102&template=example&raw=agent-loginok&raw=followme/status&num=689743&num=343"
```
