---
title: 'HTTP DELETE: HTTP 200, 202 or 204 All the Time?'
abstract: ""
created_at: 2014-06-25 12:52:00 +0000 UTC
tags:
- HTTP
- Web
slugs:
- http-delete-http-200-202-or-204-all-the-time
---

<p>I have been designing HTTP APIs (Web APIs, if you want to call it that) for a fair amount of time now and I have been handling the <a href="http://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html#sec9.7">HTTP DELETE</a> operations the same way every time. Here is a sample. <p>HTTP GET Request to get the car:</p><pre>GET http://localhost:25135/api/cars/3 HTTP/1.1
User-Agent: Fiddler
Accept: application/json
Host: localhost:25135

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 25 Jun 2014 12:36:48 GMT
Content-Length: 68

{"Id":3,"Make":"Make3","Model":"Model1","Year":2009,"Price":67437.0}</pre>
<p>HTTP DELETE Request to delete the car:<pre>DELETE http://localhost:25135/api/cars/3 HTTP/1.1
User-Agent: Fiddler
Accept: application/json
Host: localhost:25135

HTTP/1.1 204 No Content
Date: Wed, 25 Jun 2014 12:36:52 GMT</pre>
<p>Now we can see that the car is removed as I received 204 for my HTTP DELETE request. Let's send another HTTP DELETE to same resource.
<p>HTTP DELETE Request to delete the car and receive 404:<pre>DELETE http://localhost:25135/api/cars/3 HTTP/1.1
User-Agent: Fiddler
Accept: application/json
Host: localhost:25135

HTTP/1.1 404 Not Found
Date: Wed, 25 Jun 2014 12:36:52 GMT
Content-Length: 0</pre>
<p>I received 404 because /api/cars/3 is not a URI which points to a resource in my system. This is not a problem at all and it's a correct way of handling the case as I have been doing for long time now. The <a href="http://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html#sec9.1.2">idempotency</a> is also preserved because how many times you send this HTTP DELETE request, additional changes to the state of the server will not occur because the resource is already removed. So, the additional HTTP DELETE requests will just do nothing.
<p>However, here is the question in my mind: what is the real intend of the HTTP DELETE request?
<ul>
<li>Ensuring the resource is removed with the given HTTP DELETE request. 
<li>Ensuring the resource is removed.</li></ul>
<p>Here is what HTTP 1.1 spec says about HTTP DELETE:
<blockquote>
<p>The DELETE method requests that the origin server delete the resource identified by the Request-URI. This method MAY be overridden by human intervention (or other means) on the origin server. The client cannot be guaranteed that the operation has been carried out, even if the status code returned from the origin server indicates that the action has been completed successfully. However, the server SHOULD NOT indicate success unless, at the time the response is given, it intends to delete the resource or move it to an inaccessible location.
<p>A successful response SHOULD be 200 (OK) if the response includes an entity describing the status, 202 (Accepted) if the action has not yet been enacted, or 204 (No Content) if the action has been enacted but the response does not include an entity.
<p>If the request passes through a cache and the Request-URI identifies one or more currently cached entities, those entries SHOULD be treated as stale. Responses to this method are not cacheable.</p></blockquote>
<p>I don't know about you but I'm unable to figure out which two of my above intends is specified here. However, I think that the HTTP DELETE request’s intend is to ensure that the resource is removed and cannot be accessible anymore. What does this mean to my application? It means that if an HTTP DELETE operation succeeds, return a success status code (200, 202 or 204). If the resource is already removed and you receive an HTTP DELETE request for that resource, return 200 or 204 in that case; not 404. This seems more semantic to me and it is certainly be more easy for the API consumers.
<p>What do you think?
<h3>References</h3>
<ul>
<li><a href="http://stackoverflow.com/questions/4088350/is-rest-delete-really-idempotent">Is REST DELETE really idempotent?</a>
<li><a href="http://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html#sec9.1.2">Idempotent Methods</a></li></ul>  