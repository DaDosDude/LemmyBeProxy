package router

func RouteMatches(route *Route, httpMethod HttpMethod, path string) (bool, map[string]string, error) {
	if route.HttpMethod != httpMethod && route.HttpMethod != HttpMethodAll {
		return false, nil, nil
	}

	pathRegex := route.pattern
	if pathRegex == nil {
		var err error
		pathRegex, err = RegexifyRoute(route.Path)
		if err != nil {
			return false, nil, err
		}
	}

	matchesRealValues := pathRegex.FindStringSubmatch(path)
	if matchesRealValues == nil {
		return false, nil, nil
	}

	params := make(map[string]string)
	if len(matchesRealValues) > 1 {
		matchesRoute := pathRegex.FindStringSubmatch(route.Path)
		for i, value := range matchesRealValues[1:] {
			paramName := matchesRoute[i+1][1 : len(matchesRoute[i+1])-1]
			params[paramName] = value
		}
	}

	return true, params, nil
}
