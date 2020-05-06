---
title: Long-Running Asynchronous Operations, Displaying Their Events and Progress
  on Clients
abstract: I want to share a few thoughts that I have been keeping to myself on showing
  progress for long-running asyncronous operations on a system where individual events
  can be sent during ongoing operations.
created_at: 2016-07-30 16:33:00 +0000 UTC
tags:
- Architecture
- HTTP
- UX
slugs:
- long-running-asynchronous-operations-displaying-their-events-and-progress-on-clients
- long-running-asynchronous-operations-and-displaying-events-and-progress-on-clients
---

<p>With today's modern applications, we end up with lots of asynchronous operations happening and we somehow have to let the user know what is happening. This, of course, depends on the requirements and business needs. <p>Let's look at the <a href="https://www.amazon.co.uk/">Amazon</a> case. When you buy something on Amazon, Amazon tells you that your order has been taken. However, it doesn't really mean that you have actually purchased it as the process of taking the payment has not been kicked off yet. The payment process probably goes through a few stages and at the end, it has a few outcome options like success, faılure, etc. Amazon only gives the user the change to know about the outcome of this process which is done through e-mail. This is pretty straightforward to implement (famous last words?) as it doesn't require things like real-time process notifications, etc.. <p>On the other hand, creating a VM on <a href="https://azure.com">Microsoft Azure</a> has different needs. When you kick the VM creation operation through an Azure client like Azure Portal, it can give you real-time notifications on the operation status to show you in what stage the operation is (e.g. provisioning, starting up, etc.). <p>In this post, I will look at the later option and try to make a few points to prove that the implementation is a bit tricky. However, we will see a potential solution at the end and that will hopefully be helpful to you as well. <h3>Problematic Flow</h3> <p>Let's look at a problematic flow: <ul> <li>User starts an operation by sending a POST/PUT/PATCH request to an HTTP endpoint. <li>Server sends the <a href="https://tools.ietf.org/html/rfc7231#section-6.3.3">202 (Accepted)</a> response and includes a some sort of operation identifier on the response body. <li>Client subscribes to a notification channel based on this operation identifier.  <li>Client starts receiving events on that notification channel and it decides on how to display them on the user interface (UI).  <li>There are special event types which give an indication that the operation has ended (e.g. completed, failed, abandoned, etc.). </li> <ul> <li>The client needs to know about these event types and handle them specially.  <li>Whenever the client receives an event which corresponds to operation end event, it should unsubscribe from the notification channel for that specific operation.</li></ul></ul> <p>We have a few problems with the above flow: <ol> <li>Between the operation start and the client subscribing to a notification channel, there could be events that potentially could happen and the client is going to miss those. <li>If the user disconnects from the notification channel for a while (for various reasons), the client will miss the events that could happen during that disconnect period.  <li>Similar to above situation, the user might entirely shut its client and opens it up again at some point later to see the progress. In that case, we lost the events that has happened before subscribing to the notification channel on the new session.</li></ol> <h3>Possible Solution</h3> <p>The solution to all of the above problems boils down to persisting the events and opening up an HTTP endpoint to pull all the past events per each operation. Let's assume that we are performing payment transactions and want to give transparent progress on this process. Our events endpoint would look similar to below: <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>GET /payments/eefcb363/events

[
     {
          <span style="color: #a31515">"id"</span>: <span style="color: #a31515">"c12c432c"</span>,
          <span style="color: #a31515">"type"</span>: <span style="color: #a31515">"ProcessingStarted"</span>,
          <span style="color: #a31515">"message"</span>: <span style="color: #a31515">"The payment 'eefcb363' has been started for processing by worker 'h8723h7d'."</span>,
          <span style="color: #a31515">"happenedAt"</span>: <span style="color: #a31515">"2016-07-30T11:05:26.222Z"</span>
          <span style="color: #a31515">"details"</span>: {
               <span style="color: #a31515">"paymentId"</span>: <span style="color: #a31515">"eefcb363"</span>,
               <span style="color: #a31515">"workerId"</span>: <span style="color: #a31515">"h8723h7d"</span>
          }
     },

     {
          <span style="color: #a31515">"id"</span>: <span style="color: #a31515">"6bbb1d50"</span>,
          <span style="color: #a31515">"type"</span>: <span style="color: #a31515">"FraudCheckStarted"</span>,
          <span style="color: #a31515">"message"</span>: <span style="color: #a31515">"The given credit card details are being checked against fraud for payment 'eefcb363' by worker 'h8723h7d'."</span>,
          <span style="color: #a31515">"happenedAt"</span>: <span style="color: #a31515">"2016-07-30T11:05:28.779Z"</span>
          <span style="color: #a31515">"details"</span>: {
               <span style="color: #a31515">"paymentId"</span>: <span style="color: #a31515">"eefcb363"</span>,
               <span style="color: #a31515">"workerId"</span>: <span style="color: #a31515">"h8723h7d"</span>
          }
     },

     {
          <span style="color: #a31515">"id"</span>: <span style="color: #a31515">"f9e09a83"</span>,
          <span style="color: #a31515">"type"</span>: <span style="color: #a31515">"ProgressUpdated"</span>,
          <span style="color: #a31515">"message"</span>: <span style="color: #a31515">"40% of the payment 'eefcb363' has been processed by worker 'h8723h7d'."</span>,
          <span style="color: #a31515">"happenedAt"</span>: <span style="color: #a31515">"2016-07-30T11:05:29.892Z"</span>
          <span style="color: #a31515">"details"</span>: {
               <span style="color: #a31515">"paymentId"</span>: <span style="color: #a31515">"eefcb363"</span>,
               <span style="color: #a31515">"workerId"</span>: <span style="color: #a31515">"h8723h7d"</span>,
               <span style="color: #a31515">"percentage"</span>: 40
          }
     }
]</pre></div></div>
<p>This endpoint allows the client to pull the events that have already happened, and the client can mix these up with the events that it receives through the notification channel where the events are pushed in the same format.
<p>With this approach, suggested flow for the client is to follow the below steps in order to get 100% correctness on the progress and events:
<ul>
<li>User starts an operation by sending a POST/PUT/PATCH request to an HTTP endpoint.
<li>Server sends the 202 (Accepted) response and includes a some sort of operation identifier (which is 'eefcb363' in our case) on response body.
<li>Client subscribes to a notification channel based on that operation identifier.
<li>Client starts receiving events on that channel and puts them on a buffer list.
<li>Client then sends a GET request based on the operation id to get all the events which have happened so far and puts them into the buffer list (only the ones that are not already in there).
<li>Client can now start displaying the events from the buffer to the user (optionally in chronological order) and keeps receiving events through the notification channel; updates the UI simultaneously.</li></ul>
<p>Obviously, at any display stage, the client needs to honor that there are special event types which give an indication that the operation has ended. The client needs to know about these types and handle them specially (e.g. show the completed view, etc.). Besides this, whenever the client receives an event which corresponds to operation end event, it should unsubscribe from the notification channel for that specific operation.
<p>However, the situation can be a bit easier if you only care about percentage progress. If that's case, you may still want to persist all events but you only care about the latest event on the client. This makes your life easier and if you decide that your want to be more transparent about showing the progress, you can as you already have all the data for this. It is just a matter of exposing them with the above approach.
<h3>Conclusion</h3>
<p><span><span>Under any circumstances, I would suggest you to persist operation events (you can even go further with <a href="https://www.youtube.com/watch?v=JHGkaShoyNs">event sourcing</a> and make this a natural process). However, your use case may not require extensive and transparent progress reporting through the client. If that's the case, it will certainly make your implementation a lot simpler. However, you can change your mind later and go with a different progress reporting approach since you have been already persisting all the events.</span></span></p>
<p>Finally, please share your thoughts on this if you see a better way of handling this.</p>  