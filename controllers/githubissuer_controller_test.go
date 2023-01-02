package controllers

//. "github.com/onsi/ginkgo/v2"
//. "github.com/onsi/gomega"

/*var _ = Describe("GithubIssuer controller", func() {
	Context("GithubIssuer controller test", func() {

		const GithubIssuerName = "test-githubissuer"

		ctx := context.Background()
		namespace := &corev1.Namespace{}
		testCounter := 0
		typeNamespaceName := types.NamespacedName{Name: GithubIssuerName, Namespace: GithubIssuerName}

		BeforeEach(func() {
			testCounter++
			// Each test case gets its own namespace
			namespace = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: GithubIssuerName + fmt.Sprint(testCounter),
				},
			}
			err := k8sClient.Create(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))
			typeNamespaceName.Namespace = GithubIssuerName + fmt.Sprint(testCounter)
		})

		AfterEach(func() {
			err := k8sClient.Delete(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))
		})

		It("should successfully create an issue", func() {
			By("Creating the custom resource for the Kind GithubIssuer")
			githubIssuer := &githubv1.GithubIssuer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      typeNamespaceName.Name,
					Namespace: typeNamespaceName.Namespace,
				},
				Spec: githubv1.GithubIssuerSpec{
					Repo:        "https://github.com/omerbd21/transfermarkt",
					Title:       "sahar shmoroni",
					Description: "wai wai",
				},
			}
			Expect(k8sClient.Create(ctx, githubIssuer)).Should(Succeed())

		})

	})
})*/
