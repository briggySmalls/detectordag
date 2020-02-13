#!/usr/bin/env python

"""The setup script."""

from setuptools import setup, find_packages

requirements = [
    'Click>=7.0',
    'AWSIoTPythonSDK',
]

setup_requirements = ['pytest-runner', ]

test_requirements = ['pytest>=3', ]

setup(
    author="Sam Briggs",
    author_email='briggySmalls90@gmail.com',
    python_requires='>=3.5',
    classifiers=[
        'Development Status :: 2 - Pre-Alpha',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Natural Language :: English',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
        'Programming Language :: Python :: 3.8',
    ],
    description="Edge software for the detectordag powercut detector",
    entry_points={
        'console_scripts': [
            'edge=edge.cli:main',
        ],
    },
    install_requires=requirements,
    license="MIT license",
    include_package_data=True,
    keywords='edge',
    name='edge',
    packages=find_packages(include=['edge', 'edge.*']),
    setup_requires=setup_requirements,
    test_suite='tests',
    tests_require=test_requirements,
    url='https://github.com/briggySmalls/edge',
    version='0.1.0',
    zip_safe=False,
)
