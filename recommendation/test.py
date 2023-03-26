from keybert import KeyBERT

sample = '''
The Toronto Festival of Beer (TFOB), 
also known as Beer Fest, is an annual event that takes place at Exhibition Place in Toronto, Ontario, Canada. 
The festival launched in 1996 and celebrates Canada’s rich brewing history by showcasing beer of all styles, 
paired with food curated by some of Toronto's popular restaurants and Chefs, 
in addition to world renowned entertainment on the Bandshell Stage. Today, 
Toronto's Festival of Beer features more than 400 brands from around the world and many Ontario craft brewers. 
The event has become Canada’s largest beer festival with 40,000 people attending every year.
Toronto's Festival of Beer was founded in 1996 by Greg Cosway and Scott Rondeau. 
Their love for beer started at Carleton University where they started “The Gourmet Beer Club” 
which was the first of its kind in Canada. 
The festival came from those roots and has grown to become an annual celebration of the golden beverage.

In early, 2008, Greg Cosway joined forces with Les Murray, a beer industry veteran. 
The event has grown into the largest three day beer festival in Canada under the company Beerlicious. 
In 2017, Greg Cosway sold his interest in the company to pursue other business interests.
'''

kw_model = KeyBERT()
keywords = kw_model.extract_keywords(sample,stop_words=None)
print(keywords)
