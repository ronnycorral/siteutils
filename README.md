# siteutils
Tools I use to maintain my website including log parsing and creating content.
<ul>
<li>common - Golang package with functions common to multiple scripts or functions I thought would be common to multiple scripts at the time.
<li>cdlist - Creates the html for displaying a row of CD covers on site, either the 3 newest additions or 3 random covers.
<li>parsehttplogs.DEP - Parses apache logs for data used in Highcharts pages. (deprecated, replaced by 2 lambda functions)
<li>parsehttplogs - Python code for two lambda functions that (1) parses Apache logs, makes API calls to get GeoIP info and write to a file and (2) reads that data and writes it to a MySQL database. It's split in 2 because I didn't want to pay for a NAT gateway since I'm access the internet and talking to my EC2 instance.
</ul>
<p>
Coming soon:
<ul>
<li>indexsite - Indexes the site for the swish search.
<li>makesitemap - Makes the xml file used to index site by Google and Bing.
<ul>
