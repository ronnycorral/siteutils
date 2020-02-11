# siteutils
Tools I use to maintain my website including log parsing and creating content.
<ul>
<li>common - Package with functions common to multiple scripts or I thought would be common to multiple scripts at the time
<li>cdlist - Creates the html for displaying a row of CD covers on site, either the 3 newest additions or 3 random covers
<li>parsehttplogs - Parses apache logs for data used in Highcharts pages (deprecated, replaced by 2 lambda functions)
</ul>
<p>
Things to do:
<ul>
<li>indexsite - Indexes the site for the swish search
<li>makesitemap - makes the xml files used to index site by Google and Bing.
<li>The 2 lambda functions that (1) parses the log data, makes API calls to get the GeoIP info and write to a file and (2) reads that data and writes to to a MySQL database
<ul>

