---
title: Elasticsearch Array Contains Search With Terms Filter
abstract: Here is a quick blog post on Elasticsearch and terms filter to achieve array
  contains search with terms filter
created_at: 2015-07-14 22:17:00 +0000 UTC
tags:
- Elasticsearch
slugs:
- elasticsearch-array-contains-search-with-terms-filter
---

<p>Here is a quick blog post on <a href="https://www.elastic.co/products/elasticsearch">Elasticsearch</a> and <a href="https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-terms-filter.html">terms filter</a> while I still remember how the hell it works :) Yes, this is possibly the 20th time that I looked for how to achieve array contains functionality in Elasticseach and it's a clear sign for me that I need to blog about it :)</p> <p>I created the index called movies (mostly borrowed from <a href="http://joelabrahamsson.com/elasticsearch-101/">Joel's great Elasticsearch 101 blog post</a>) and here is its mapping:</p> <div class="code-wrapper border-shadow-1"> <div style="color: black; background-color: white"><pre>PUT movies/_mapping/movie
{
  <span style="color: #a31515">"movie"</span>: {
    <span style="color: #a31515">"properties"</span>: {
       <span style="color: #a31515">"director"</span>: {
          <span style="color: #a31515">"type"</span>: <span style="color: #a31515">"string"</span>
       },
       <span style="color: #a31515">"genres"</span>: {
          <span style="color: #a31515">"type"</span>: <span style="color: #a31515">"string"</span>,
          <span style="color: #a31515">"index"</span>: <span style="color: #a31515">"not_analyzed"</span>
       },
       <span style="color: #a31515">"title"</span>: {
          <span style="color: #a31515">"type"</span>: <span style="color: #a31515">"string"</span>
       },
       <span style="color: #a31515">"year"</span>: {
          <span style="color: #a31515">"type"</span>: <span style="color: #a31515">"long"</span>
       }
    }
  }
}</pre></div></div>
<p>The genres field mapping is important here as it needs to be <a href="https://www.elastic.co/guide/en/elasticsearch/guide/current/mapping-intro.html#_index_2">not analyzed</a>. I also indexed a few stuff in it:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>POST movies/movie
{
    <span style="color: #a31515">"title"</span>: <span style="color: #a31515">"Apocalypse Now"</span>,
    <span style="color: #a31515">"director"</span>: <span style="color: #a31515">"Francis Ford Coppola"</span>,
    <span style="color: #a31515">"year"</span>: 1979,
    <span style="color: #a31515">"genres"</span>: [<span style="color: #a31515">"Drama"</span>, <span style="color: #a31515">"War"</span>, <span style="color: #a31515">"Foo"</span>]
}

POST movies/movie
{
    <span style="color: #a31515">"title"</span>: <span style="color: #a31515">"Apocalypse Now"</span>,
    <span style="color: #a31515">"director"</span>: <span style="color: #a31515">"Francis Ford Coppola"</span>,
    <span style="color: #a31515">"year"</span>: 1979,
    <span style="color: #a31515">"genres"</span>: [<span style="color: #a31515">"Drama"</span>, <span style="color: #a31515">"War"</span>, <span style="color: #a31515">"Foo"</span>, <span style="color: #a31515">"Bar"</span>]
}

POST movies/movie
{
    <span style="color: #a31515">"title"</span>: <span style="color: #a31515">"Apocalypse Now"</span>,
    <span style="color: #a31515">"director"</span>: <span style="color: #a31515">"Francis Ford Coppola"</span>,
    <span style="color: #a31515">"year"</span>: 1979,
    <span style="color: #a31515">"genres"</span>: [<span style="color: #a31515">"Drama"</span>, <span style="color: #a31515">"Bar"</span>]
}</pre></div></div>
<p>Now, I am interested in finding out the movies which is in <em><strong>War </strong></em>or <em><strong>Foo</strong></em> genre. The way to achieve that is the terms filter as mentioned:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>GET movies/movie/_search
{
  <span style="color: #a31515">"query"</span>: {
    <span style="color: #a31515">"filtered"</span>: {
      <span style="color: #a31515">"query"</span>: {
        <span style="color: #a31515">"match_all"</span>: {}
      },
      <span style="color: #a31515">"filter"</span>: {
        <span style="color: #a31515">"terms"</span>: {
          <span style="color: #a31515">"genres"</span>: [<span style="color: #a31515">"War"</span>, <span style="color: #a31515">"Foo"</span>]
        }
      }
    }
  }
}</pre></div></div>
<p>We will get us the following result:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: #a31515">"hits"</span>: [
   {
      <span style="color: #a31515">"_index"</span>: <span style="color: #a31515">"movies"</span>,
      <span style="color: #a31515">"_type"</span>: <span style="color: #a31515">"movie"</span>,
      <span style="color: #a31515">"_id"</span>: <span style="color: #a31515">"AU6OkygJidzUtyfB9L2D"</span>,
      <span style="color: #a31515">"_score"</span>: 1,
      <span style="color: #a31515">"_source"</span>: {
         <span style="color: #a31515">"title"</span>: <span style="color: #a31515">"Apocalypse Now"</span>,
         <span style="color: #a31515">"director"</span>: <span style="color: #a31515">"Francis Ford Coppola"</span>,
         <span style="color: #a31515">"year"</span>: 1979,
         <span style="color: #a31515">"genres"</span>: [
            <span style="color: #a31515">"Drama"</span>,
            <span style="color: #a31515">"War"</span>,
            <span style="color: #a31515">"Foo"</span>,
            <span style="color: #a31515">"Bar"</span>
         ]
      }
   },
   {
      <span style="color: #a31515">"_index"</span>: <span style="color: #a31515">"movies"</span>,
      <span style="color: #a31515">"_type"</span>: <span style="color: #a31515">"movie"</span>,
      <span style="color: #a31515">"_id"</span>: <span style="color: #a31515">"AU6OkwUeidzUtyfB9L1q"</span>,
      <span style="color: #a31515">"_score"</span>: 1,
      <span style="color: #a31515">"_source"</span>: {
         <span style="color: #a31515">"title"</span>: <span style="color: #a31515">"Apocalypse Now"</span>,
         <span style="color: #a31515">"director"</span>: <span style="color: #a31515">"Francis Ford Coppola"</span>,
         <span style="color: #a31515">"year"</span>: 1979,
         <span style="color: #a31515">"genres"</span>: [
            <span style="color: #a31515">"Drama"</span>,
            <span style="color: #a31515">"War"</span>,
            <span style="color: #a31515">"Foo"</span>
         ]
      }
   }
]</pre></div></div>
<p>What if we want to see the movies which is in <strong><em>War,</em></strong> <strong><em>Foo </em></strong>and<strong><em> Bar</em></strong> genres at the same time? Well, there are probably other ways of doing this but here is how I hacked it together with <a href="https://www.elastic.co/guide/en/elasticsearch/reference/1.6/query-dsl-bool-filter.html">bool</a> and <a href="https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-term-filter.html">term filter</a>:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre>GET movies/movie/_search
{
  <span style="color: #a31515">"query"</span>: {
    <span style="color: #a31515">"filtered"</span>: {
      <span style="color: #a31515">"query"</span>: {
        <span style="color: #a31515">"match_all"</span>: {}
      },
      <span style="color: #a31515">"filter"</span>: {
        <span style="color: #a31515">"bool"</span>: {
          <span style="color: #a31515">"must"</span>: [
            { 
              <span style="color: #a31515">"term"</span>: {
                <span style="color: #a31515">"genres"</span>: <span style="color: #a31515">"War"</span>
              }
            }
          ],
          
          <span style="color: #a31515">"must"</span>: [
            { 
              <span style="color: #a31515">"term"</span>: {
                <span style="color: #a31515">"genres"</span>: <span style="color: #a31515">"Foo"</span>
              }
            }
          ],
          
          <span style="color: #a31515">"must"</span>: [
            { 
              <span style="color: #a31515">"term"</span>: {
                <span style="color: #a31515">"genres"</span>: <span style="color: #a31515">"Bar"</span>
              }
            }
          ]
        }
      }
    }
  }
}</pre></div></div>
<p>The result:</p>
<div class="code-wrapper border-shadow-1">
<div style="color: black; background-color: white"><pre><span style="color: #a31515">"hits"</span>: [
   {
      <span style="color: #a31515">"_index"</span>: <span style="color: #a31515">"movies"</span>,
      <span style="color: #a31515">"_type"</span>: <span style="color: #a31515">"movie"</span>,
      <span style="color: #a31515">"_id"</span>: <span style="color: #a31515">"AU6OkygJidzUtyfB9L2D"</span>,
      <span style="color: #a31515">"_score"</span>: 1,
      <span style="color: #a31515">"_source"</span>: {
         <span style="color: #a31515">"title"</span>: <span style="color: #a31515">"Apocalypse Now"</span>,
         <span style="color: #a31515">"director"</span>: <span style="color: #a31515">"Francis Ford Coppola"</span>,
         <span style="color: #a31515">"year"</span>: 1979,
         <span style="color: #a31515">"genres"</span>: [
            <span style="color: #a31515">"Drama"</span>,
            <span style="color: #a31515">"War"</span>,
            <span style="color: #a31515">"Foo"</span>,
            <span style="color: #a31515">"Bar"</span>
         ]
      }
   }
]</pre></div></div>
<p><a href="https://www.elastic.co/guide/en/elasticsearch/guide/current/_finding_multiple_exact_values.html">The exact match is a whole different story</a> :)</p>  