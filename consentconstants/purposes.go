package consentconstants

// Purpose is one of the IAB GDPR purposes. These appear in:
//   1. `root.purposes[i]` of the vendor list: https://vendorlist.consensu.org/vendorlist.json
//   2. PurposesAllowed of the Consent string: https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/master/Consent%20string%20and%20vendor%20list%20formats%20v1.1%20Final.md#vendor-consent-string-format-
type Purpose uint8

// TCF 1 Purposes:
const (
	// InfoStorageAccess includes the storage of information, or access to information that is already stored,
	// on your device such as advertising identifiers, device identifiers, cookies, and similar technologies.
	InfoStorageAccess Purpose = 1

	// Personalization includes the collection and processing of information about your use of this service to subsequently personalise
	// advertising and/or content for you in other contexts, such as on other websites or apps, over time.
	// Typically, the content of the site or app is used to make inferences about your interests, which inform
	// future selection of advertising and/or content.
	Personalization Purpose = 2

	// AdSelectionDeliveryReporting includes the collection of information, and combination with previously collected information, to select and deliver
	// advertisements for you, and to measure the delivery and effectiveness of such advertisements. This includes using previously
	// collected information about your interests to select ads, processing data about what advertisements were shown, how often they were shown,
	// when and where they were shown, and whether you took any action related to the advertisement, including for example clicking an ad or making a purchase.
	// This does not include personalisation, which is the collection and processing of information about your use of this service to
	// subsequently personalise advertising and/or content for you in other contexts, such as websites or apps, over time.
	AdSelectionDeliveryReporting Purpose = 3

	// ContentSelectionDeliveryReporting includes the collection of information, and combination with previously collected information, to select and deliver
	// content for you, and to measure the delivery and effectiveness of such content. This includes using previously
	// collected information about your interests to select content, processing data about what content was shown,
	// how often or how long it was shown, when and where it was shown, and whether the you took any action related to
	// the content, including for example clicking on content. This does not include personalisation, which is the collection
	// and processing of information about your use of this service to subsequently personalise content and/or advertising
	// for you in other contexts, such as websites or apps, over time.
	ContentSelectionDeliveryReporting Purpose = 4

	// Measurement includes the collection of information about your use of the content, and combination with previously collected information,
	// used to measure, understand, and report on your usage of the service. This does not include personalisation, the
	// collection of information about your use of this service to subsequently personalise content and/or advertising
	// for you in other contexts, i.e. on other service, such as websites or apps, over time.
	Measurement Purpose = 5
)

// TCF 2.0 Purposes:
const (
	// Cookies, device identifiers, or other information can be stored or accessed on your device for the purposes presented to you.
	// Vendors can:
	//   * Store and access information on the device such as cookies and device identifiers presented to a user.
	// Reuse InfoStorageAccess above

	// Ads can be shown to you based on the content you are viewing, the app you are using, your approximate location, or your device type.
	// To do basic ad selection vendors can:
	//   * Use real-time information about the context in which the ad will be shown, to show the ad, including information about the content and
	//     the device, such as: device type and capabilities, user agent, URL, IP address
	//   * Use a user's non-precise geolocation data
	//   * Control the frequency of ads shown to a user.\n* Sequence the order in which ads are shown to a user.
	//   * Prevent an ad from serving in an unsuitable editorial (brand-unsafe) context
	// Vendors cannot:
	//   * Create a personalised ads profile using this information for the selection of future ads.
	//   * N.B. Non-precise means only an approximate location involving at least a radius of 500 meters is permitted.
	BasicAdserving Purpose = 2

	// A profile can be built about you and your interests to show you personalised ads that are relevant to you.
	// To create a personalised ads profile vendors can:
	//   * Collect information about a user, including a user's activity, interests, demographic information, or location, to create or edit a user profile for use in personalised advertising.
	//   * Combine this information with other information previously collected, including from across websites and apps, to create or edit a user profile for use in personalised advertising.
	PersonalizationProfile Purpose = 3

	// Personalised ads can be shown to you based on a profile about you.
	// To select personalised ads vendors can:
	//   * Select personalised ads based on a user profile or other historical user data, including a user's prior activity, interests, visits to sites or apps, location, or demographic information.
	PersonalizationSelection Purpose = 4

	// A profile can be built about you and your interests to show you personalised content that is relevant to you.
	// To create a personalised content profile vendors can:
	//   * Collect information about a user, including a user's activity, interests, visits to sites or apps, demographic information, or location, to create or edit a user profile for personalising content.
	//   * Combine this information with other information previously collected, including from across websites and apps, to create or edit a user profile for use in personalising content.
	ContentProfile Purpose = 5

	// Personalised content can be shown to you based on a profile about you.
	// To select personalised content vendors can:
	//   * Select personalised content based on a user profile or other historical user data, including a user\u2019s prior activity, interests, visits to sites or apps, location, or demographic information.
	ContentSelection Purpose = 6

	// The performance and effectiveness of ads that you see or interact with can be measured.
	// To measure ad performance vendors can:
	//   * Measure whether and how ads were delivered to and interacted with by a user
	//   * Provide reporting about ads including their effectiveness and performance
	//   * Provide reporting about users who interacted with ads using data observed during the course of the user's interaction with that ad
	//   * Provide reporting to publishers about the ads displayed on their property
	//   * Measure whether an ad is serving in a suitable editorial environment (brand-safe) context
	//   * Determine the percentage of the ad that had the opportunity to be seen and the duration of that opportunity
	//   * Combine this information with other information previously collected, including from across websites and apps
	// Vendors cannot:
	//   *Apply panel- or similarly-derived audience insights data to ad measurement data without a Legal Basis to apply market research to generate audience insights (Purpose 9)
	AdPerformance Purpose = 7

	// The performance and effectiveness of content that you see or interact with can be measured.
	// To measure content performance vendors can:
	//    * Measure and report on how content was delivered to and interacted with by users.
	//    * Provide reporting, using directly measurable or known information, about users who interacted with the content
	//    * Combine this information with other information previously collected, including from across websites and apps.
	// Vendors cannot:
	//    * Measure whether and how ads (including native ads) were delivered to and interacted with by a user.
	//    * Apply panel- or similarly derived audience insights data to ad measurement data without a Legal Basis to apply market research to generate audience insights (Purpose 9)
	ContentPerformance Purpose = 8

	// Market research can be used to learn more about the audiences who visit sites/apps and view ads.
	// To apply market research to generate audience insights vendors can:
	//    * Provide aggregate reporting to advertisers or their representatives about the audiences reached by their ads, through panel-based and similarly derived insights.
	//    * Provide aggregate reporting to publishers about the audiences that were served or interacted with content and/or ads on their property by applying panel-based and similarly derived insights.
	//    * Associate offline data with an online user for the purposes of market research to generate audience insights if vendors have declared to match and combine offline data sources (Feature 1)
	//    * Combine this information with other information previously collected including from across websites and apps.
	// Vendors cannot:
	//    * Measure the performance and effectiveness of ads that a specific user was served or interacted with, without a Legal Basis to measure ad performance.
	//    * Measure which content a specific user was served and how they interacted with it, without a Legal Basis to measure content performance.
	MarketResearch Purpose = 9

	// Your data can be used to improve existing systems and software, and to develop new products
	// To develop new products and improve products vendors can:
	//    * Use information to improve their existing products with new features and to develop new products
	//    * Create new models and algorithms through machine learning
	// Vendors cannot:
	//    * Conduct any other data processing operation allowed under a different purpose under this purpose
	DevelopImprove Purpose = 10
)
