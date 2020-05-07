# siteutils
Tools I use to maintain my website including log parsing and creating content.
<ul>
<li>common - Golang package with functions common to multiple scripts or functions I thought would be common to multiple scripts at the time.
<li>cdlist - Creates the html for displaying a row of CD covers on site, either the 3 newest additions or 3 random covers.
<li>parsehttplogs.DEPRECATED - Parses apache logs to create data used to populate the Highcharts pages. (deprecated, replaced by 2 lambda functions)
<li>parsehttplogs - Python code for two lambda functions that (1) parses Apache logs, makes API calls to get GeoIP info and writes that data to a file and (2) reads those files and writes the data to a MySQL database. It's split in 2 because I didn't want to pay for a NAT gateway.
</ul>
<p>
Coming soon:
<ul>
<li>indexsite - Indexes the site for the swish search.
<li>makesitemap - Makes the xml files used to index site by Google and Bing.
<ul>
